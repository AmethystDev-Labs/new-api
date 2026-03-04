package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/QuantumNous/new-api/constant"
)

// ValidateRedirectURL validates that a redirect URL is safe to use.
// It checks that:
//   - The URL is properly formatted
//   - The scheme is either http or https
//   - The domain is in the trusted domains list (exact match or subdomain)
//
// Returns nil if the URL is valid and trusted, otherwise returns an error
// describing why the validation failed.
func ValidateRedirectURL(rawURL string) error {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %s", err.Error())
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: only http and https are allowed")
	}

	domain := strings.ToLower(parsedURL.Hostname())

	for _, trustedDomain := range constant.TrustedRedirectDomains {
		if domainMatchesTrustedDomain(domain, trustedDomain) {
			return nil
		}
	}

	return fmt.Errorf("domain %s is not in the trusted domains list", domain)
}

// ValidateAMFSBaseURL validates AMFS API base URL and returns normalized base URL.
// Security rules:
//   - must be valid absolute HTTPS URL
//   - must not contain query, fragment, or userinfo
//   - host must be in TrustedAMFSDomains (exact or subdomain match)
func ValidateAMFSBaseURL(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", fmt.Errorf("AMFS API Base is empty")
	}

	parsedURL, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid AMFS API Base format: %s", err.Error())
	}
	if parsedURL.Scheme != "https" {
		return "", fmt.Errorf("invalid AMFS API Base scheme: only https is allowed")
	}
	if parsedURL.Hostname() == "" {
		return "", fmt.Errorf("invalid AMFS API Base: hostname is required")
	}
	if parsedURL.User != nil {
		return "", fmt.Errorf("invalid AMFS API Base: userinfo is not allowed")
	}
	if parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return "", fmt.Errorf("invalid AMFS API Base: query and fragment are not allowed")
	}
	if parsedURL.Path != "" && parsedURL.Path != "/" {
		return "", fmt.Errorf("invalid AMFS API Base: path is not allowed")
	}
	if parsedURL.Port() != "" && parsedURL.Port() != "443" {
		return "", fmt.Errorf("invalid AMFS API Base port: only 443 is allowed")
	}

	domain := strings.ToLower(parsedURL.Hostname())
	for _, trustedDomain := range constant.TrustedAMFSDomains {
		if domainMatchesTrustedDomain(domain, trustedDomain) {
			return strings.TrimRight(parsedURL.String(), "/"), nil
		}
	}
	return "", fmt.Errorf("AMFS API Base domain %s is not in the trusted domains list", domain)
}

func domainMatchesTrustedDomain(domain string, trustedDomain string) bool {
	td := strings.ToLower(strings.TrimSpace(trustedDomain))
	if td == "" {
		return false
	}
	return domain == td || strings.HasSuffix(domain, "."+td)
}
