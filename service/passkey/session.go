package passkey

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	webauthn "github.com/go-webauthn/webauthn/webauthn"
)

var errSessionNotFound = errors.New("Passkey 会话不存在或已过期")

func SaveSessionData(c *gin.Context, key string, data *webauthn.SessionData) error {
	session := sessions.Default(c)
	if data == nil {
		session.Delete(key)
		return session.Save()
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	session.Set(key, string(payload))
	return session.Save()
}

func PopSessionData(c *gin.Context, key string) (*webauthn.SessionData, error) {
	session := sessions.Default(c)
	raw := session.Get(key)
	if raw == nil {
		return nil, errSessionNotFound
	}
	session.Delete(key)
	if err := session.Save(); err != nil {
		return nil, fmt.Errorf("failed to save session after deleting invalid passkey state: %w", err)
	}
	var data webauthn.SessionData
	switch value := raw.(type) {
	case string:
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			return nil, err
		}
	case []byte:
		if err := json.Unmarshal(value, &data); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Passkey 会话格式无效")
	}
	return &data, nil
}
