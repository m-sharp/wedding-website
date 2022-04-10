package lib

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
)

const (
	WeddingAddress = "lindenandmike@gmail.com"

	GMailHost = "smtp.gmail.com"
	GMailPort = 587

	passwordEnvVar = "EMAILPASSWORD"
)

type EmailFailedError struct {
	inner error
}

func (e *EmailFailedError) Error() string {
	return fmt.Sprintf("failed to send email over smtp: %s", e.inner)
}

func (e *EmailFailedError) Unwrap() error {
	return e.inner
}

func SendEmail(message string, to ...string) error {
	if len(to) < 1 {
		return &EmailFailedError{inner: errors.New("no To addresses specified")}
	}

	password, err := getPass()
	if err != nil {
		return &EmailFailedError{inner: err}
	}
	if err = smtp.SendMail(
		fmt.Sprintf("%s:%v", GMailHost, GMailPort),
		smtp.PlainAuth("", WeddingAddress, password, GMailHost),
		WeddingAddress,
		to,
		[]byte(message),
	); err != nil {
		return &EmailFailedError{inner: err}
	}
	return nil
}

func getPass() (string, error) {
	value, ok := os.LookupEnv(passwordEnvVar)
	if !ok {
		return "", fmt.Errorf("ENVVAR for %q not found", passwordEnvVar)
	}
	return value, nil
}
