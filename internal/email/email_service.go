package email

import (
	"bytes"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
)

type EmailService interface {
	SendTemplatedEmail(toEmail, subject, templateName string, variables map[string]string) error
}

type EmailServiceImpl struct {
	fromEmail string
	dialer    *gomail.Dialer
}

func NewEmailService() *EmailServiceImpl {
	fromEmail := os.Getenv("EMAIL_FROM")
	host := os.Getenv("SMTP_HOST")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		logger.Logger.Error("Invalid SMTP port", "error", err)
		panic(1)
	}

	dialer := gomail.NewDialer(host, port, username, password)
	return &EmailServiceImpl{fromEmail: fromEmail, dialer: dialer}
}

func (s *EmailServiceImpl) SendTemplatedEmail(toEmail, subject, templateName string, variables map[string]string) error {
	logger.Logger.Info("Sending email", "to", toEmail)

	body, err := s.buildEmailBody(templateName, variables)
	if err != nil {
		logger.Logger.Error("Failed to build email body", "to", toEmail, "error", err)
		return service_errors.NewErrInternalServer("Failed to build email body")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		logger.Logger.Error("Failed to send email", "to", toEmail, "error", err)
		return service_errors.NewErrInternalServer("Failed to send email")
	}

	return nil
}

func (s *EmailServiceImpl) buildEmailBody(templateName string, data map[string]string) (string, error) {
	filePath := "templates/" + templateName

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		logger.Logger.Error("Failed to parse template", "template", templateName, "error", err.Error())
		return "", service_errors.NewErrInternalServer("Failed to parse email template")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		logger.Logger.Error("Failed to execute template", "template", templateName, "error", err.Error())
		return "", service_errors.NewErrInternalServer("Failed to execute email")
	}

	return buf.String(), nil
}
