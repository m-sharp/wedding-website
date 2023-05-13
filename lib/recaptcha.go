package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

const (
	recaptchaUrl = "https://www.google.com/recaptcha/api/siteverify"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func Verify(log *zap.Logger, recaptchaSecret, userRespToken, userIp string) (bool, error) {
	bodyData := url.Values{}
	bodyData.Set("secret", recaptchaSecret)
	bodyData.Set("response", userRespToken)
	bodyData.Set("remoteip", userIp)

	resp, err := http.Post(
		recaptchaUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(bodyData.Encode()),
	)
	if err != nil {
		return false, fmt.Errorf("error sending recaptcha request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad response code returned by recaptcha request: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read recaptcha response body: %w", err)
	}

	var recaptchaResp RecaptchaResponse
	if err := json.Unmarshal(body, &recaptchaResp); err != nil {
		return false, fmt.Errorf("failed to decode recaptcha response: %w", err)
	}

	log.Info("Got recaptcha response", zap.Any("Response", recaptchaResp))
	return recaptchaResp.Success, nil
}
