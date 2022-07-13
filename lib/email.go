package lib

import (
	"errors"
	"fmt"
	"net/smtp"
	"strconv"

	"go.uber.org/zap"
)

const (
	WeddingAddress = "lindenandmike@gmail.com"

	GMailHost = "smtp.gmail.com"
	GMailPort = 587
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

func SendEmail(cfg *Config, log *zap.Logger, message string, to ...string) error {
	if len(to) < 1 {
		return &EmailFailedError{inner: errors.New("no To addresses specified")}
	}

	password, err := cfg.Get(EmailPass)
	if err != nil {
		return &EmailFailedError{inner: err}
	}

	log = log.Named("Emailer")
	log.Info(
		"Sending Email",
		zap.String("From", WeddingAddress),
		zap.String("SMTP Host", GMailHost),
		zap.String("SMTP Port", strconv.Itoa(GMailPort)),
		zap.Any("To", to),
	)
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
