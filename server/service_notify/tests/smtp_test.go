package tests

import (
	"context"
	"net/smtp"
	"testing"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
	implApplication "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service/impl"
	infraMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/mail"
	domainMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mail"
)

/**
 * Test send mail SMTP
 */
func TestSendMailSMTP(t *testing.T) {
	userName := "lytranvinh.work@gmail.com"
	password := "" // replace with your app password
	host := "smtp.gmail.com"
	port := 587
	to := []string{"nguyenducquanghihihi@gmail.com", "udontknowmebaby113@gmail.com"}
	smtpAuth := smtp.PlainAuth(
		"",
		userName,
		password,
		host,
	)
	service := infraMail.NewSMTPMail(
		smtpAuth,
		host, 
		port,
	)
	htmlContent := `<h1>This is a test email</h1><p>Sent using Go's net/smtp package.</p>`
	subject := "Test Email from Go"
	if err := service.SendMail(to, subject, htmlContent); err != nil {
		t.Errorf("Failed to send email: %v", err)
	} else {
		t.Log("Email sent successfully")
	}
}

func TestSendMailSMTPApplication(t *testing.T) {
	userName := "lytranvinh.work@gmail.com"
	password := "" // replace with your app password
	host := "smtp.gmail.com"
	port := 587
	to := []string{"udontknowmebaby113@gmail.com"}
	// init domain mail
	smtpAuth := smtp.PlainAuth(
		"",
		userName,
		password,
		host,
	)
	if err := domainMail.SetSMTPService(infraMail.NewSMTPMail(
		smtpAuth,
		host,
		port,
	)); err != nil {
		t.Errorf("Failed to set SMTP service: %v", err)
	}
	if err := domainMail.SetHtmlMailContent(infraMail.NewHTMLContentMail()); err != nil {
		t.Errorf("Failed to set HTML mail content service: %v", err)
	}
	// test send mail
	impl := implApplication.NewMailService()
	ctx := context.Background()
	err := impl.SendMessageForgotPassword(
		ctx,
		applicationModel.MailForgotPassword{
			To:          to[0],
			UrlAuth:     "https://example.com/reset-password",
			NewPassword: "new_secure_password",
			Expired:     10,
		},
	)
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
	} else {
		t.Log("Email sent successfully")
	}
}
