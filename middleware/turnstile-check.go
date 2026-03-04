package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type turnstileCheckResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

type amfsScoreRequest struct {
	EventID string `json:"eventId"`
	SiteID  string `json:"siteId"`
	Scene   string `json:"scene"`
	UserID  string `json:"userId,omitempty"`
}

type amfsScoreResponse struct {
	Score float64 `json:"score"`
}

func TurnstileCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !common.TurnstileCheckEnabled {
			c.Next()
			return
		}

		session := sessions.Default(c)
		captchaChecked := session.Get("captcha")
		legacyTurnstileChecked := session.Get("turnstile")
		if captchaChecked != nil || legacyTurnstileChecked != nil {
			c.Next()
			return
		}

		provider := resolveCaptchaProvider()
		if provider == "amfs" && c.FullPath() != "/api/user/login" {
			c.Next()
			return
		}

		token := getCaptchaTokenFromQuery(c, provider)
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": getCaptchaEmptyMessage(provider),
			})
			c.Abort()
			return
		}

		if provider == "amfs" {
			blocked, err := verifyAMFSByScore(c, token)
			if blocked {
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": "验证码校验失败，请刷新重试！",
				})
				c.Abort()
				return
			}
			if err != nil {
				common.SysLog("amfs check failed: " + err.Error())
				if !common.AMFSFailOpen {
					c.JSON(http.StatusOK, gin.H{
						"success": false,
						"message": "AMFS 校验失败，请刷新重试！",
					})
					c.Abort()
					return
				}
			}
		} else {
			verifyErr := verifyTurnstileOrHCaptcha(c, provider, token)
			if verifyErr != nil {
				common.SysLog(verifyErr.Error())
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": getCaptchaFailedMessage(provider),
				})
				c.Abort()
				return
			}
		}

		session.Set("captcha", true)
		err := session.Save()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "无法保存会话信息，请重试",
				"success": false,
			})
			return
		}
		c.Next()
	}
}

func resolveCaptchaProvider() string {
	provider := common.CaptchaProvider
	if provider != "hcaptcha" && provider != "amfs" {
		return "turnstile"
	}
	return provider
}

func getCaptchaTokenFromQuery(c *gin.Context, provider string) string {
	token := c.Query("captcha")
	if token != "" {
		return token
	}

	switch provider {
	case "hcaptcha":
		token = c.Query("hcaptcha")
		if token == "" {
			token = c.Query("turnstile")
		}
	case "amfs":
		token = c.Query("amfs")
		if token == "" {
			token = c.Query("eventId")
		}
	default:
		token = c.Query("turnstile")
		if token == "" {
			token = c.Query("hcaptcha")
		}
	}
	return token
}

func getCaptchaEmptyMessage(provider string) string {
	switch provider {
	case "hcaptcha":
		return "hCaptcha token 为空"
	case "amfs":
		return "AMFS eventId 为空"
	default:
		return "Turnstile token 为空"
	}
}

func getCaptchaFailedMessage(provider string) string {
	if provider == "hcaptcha" {
		return "hCaptcha 校验失败，请刷新重试！"
	}
	return "Turnstile 校验失败，请刷新重试！"
}

func verifyTurnstileOrHCaptcha(c *gin.Context, provider string, token string) error {
	verifyURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	secret := common.TurnstileSecretKey
	if provider == "hcaptcha" {
		verifyURL = "https://hcaptcha.com/siteverify"
		secret = common.HCaptchaSecretKey
	}
	if strings.TrimSpace(secret) == "" {
		return fmt.Errorf("%s secret key is empty", provider)
	}

	rawRes, err := http.PostForm(verifyURL, url.Values{
		"secret":   {secret},
		"response": {token},
		"remoteip": {c.ClientIP()},
	})
	if err != nil {
		return err
	}
	defer rawRes.Body.Close()

	var res turnstileCheckResponse
	if err = common.DecodeJson(rawRes.Body, &res); err != nil {
		return err
	}
	if !res.Success {
		if len(res.ErrorCodes) > 0 {
			return fmt.Errorf("%s verification failed: %s", provider, strings.Join(res.ErrorCodes, ","))
		}
		return fmt.Errorf("%s verification failed", provider)
	}
	return nil
}

func verifyAMFSByScore(c *gin.Context, eventID string) (bool, error) {
	if strings.TrimSpace(common.AMFSApiBase) == "" || strings.TrimSpace(common.AMFSSiteID) == "" {
		return false, fmt.Errorf("amfs config missing")
	}
	apiBase, validateErr := common.ValidateAMFSBaseURL(common.AMFSApiBase)
	if validateErr != nil {
		return false, validateErr
	}

	timeout := time.Duration(common.AMFSTimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	client := &http.Client{Timeout: timeout}
	scoreReq := amfsScoreRequest{
		EventID: eventID,
		SiteID:  common.AMFSSiteID,
		Scene:   "login",
	}

	reqBody, err := common.Marshal(scoreReq)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, apiBase+"/v1/score", bytes.NewReader(reqBody))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	rawRes, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer rawRes.Body.Close()

	if rawRes.StatusCode != http.StatusOK {
		return false, fmt.Errorf("amfs score status code %d", rawRes.StatusCode)
	}

	var scoreRes amfsScoreResponse
	if err = common.DecodeJson(rawRes.Body, &scoreRes); err != nil {
		return false, err
	}

	if scoreRes.Score > float64(common.AMFSScoreThreshold) {
		common.SysLog("amfs blocked by score(" + strconv.FormatFloat(scoreRes.Score, 'f', -1, 64) + ") > threshold(" + strconv.Itoa(common.AMFSScoreThreshold) + ")")
		return true, nil
	}
	return false, nil
}
