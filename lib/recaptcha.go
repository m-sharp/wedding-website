package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
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
	bodyData := struct {
		Secret   string `json:"secret"`
		Response string `json:"response"`
		RemoteIP string `json:"remoteip"`
	}{
		Secret:   recaptchaSecret,
		Response: userRespToken,
		RemoteIP: userIp,
	}
	content, err := json.Marshal(bodyData)
	if err != nil {
		return false, fmt.Errorf("failed to construct recaptcha body data: %w", err)
	}

	resp, err := http.Post(recaptchaUrl, "application/json", bytes.NewReader(content))
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
