package email

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"go.uber.org/zap"

	"github.com/m-sharp/wedding-website/lib"
)

const (
	WeddingAddress = "lindenandmike@gmail.com" // ToDo: Configurable?

	GMailHost = "smtp.gmail.com"
	GMailPort = 587
)

type MessageTemplate string

const (
	NotificationTmpl = MessageTemplate("notification")
	ResponseTmpl     = MessageTemplate("response")
)

var (
	emailTemplatesDir = filepath.FromSlash(filepath.Join("lib", "email"))
	templatePaths     = map[MessageTemplate]string{
		NotificationTmpl: filepath.FromSlash(filepath.Join(emailTemplatesDir, "emailNotification.tmpl")),
		ResponseTmpl:     filepath.FromSlash(filepath.Join(emailTemplatesDir, "emailResponse.tmpl")),
	}
)

func SendEmail(
	cfg *lib.Config,
	log *zap.Logger,
	toAddress string,
	targetTemplate MessageTemplate,
	renderData map[string]interface{},
) error {
	log = log.Named("Emailer")

	if toAddress == "" {
		return errors.New("no To address specified")
	}

	password, err := cfg.Get(lib.EmailPass)
	if err != nil {
		return fmt.Errorf("failed to get email account password from config: %w", err)
	}

	targetPath, ok := templatePaths[targetTemplate]
	if !ok {
		return errors.New("target email template not found")
	}

	tmpl, err := template.ParseFiles(targetPath)
	if err != nil {
		return fmt.Errorf("failed to parse notification email template: %w", err)
	}

	renderData["From"] = WeddingAddress
	renderData["To"] = toAddress

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(
		buf,
		string(targetTemplate),
		renderData,
	); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}
	content := buf.String()
	content = strings.TrimSpace(content)

	log.Info(
		"Sending Email",
		zap.String("From", WeddingAddress),
		zap.String("SMTP Host", GMailHost),
		zap.String("SMTP Port", strconv.Itoa(GMailPort)),
		zap.String("To", toAddress),
	)
	if err = smtp.SendMail(
		fmt.Sprintf("%s:%v", GMailHost, GMailPort),
		smtp.PlainAuth("", WeddingAddress, password, GMailHost),
		WeddingAddress,
		[]string{toAddress},
		[]byte(content),
	); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
