package impl

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service"
	domainMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mail"
)

// Structure mail service
type MailService struct {
}

// SendMessageForgotPassword implements service.IMailService.
func (m *MailService) SendMessageForgotPassword(ctx context.Context, input model.MailForgotPassword) error {
	htmlGen := domainMail.GetHtmlMailContent()
	htmlContent, err := htmlGen.ForgotPassword(
		input.To,
		input.UrlAuth,
		input.NewPassword,
		input.Expired,
	)
	if err != nil {
		return err
	}
	smtpService := domainMail.GetSMTPService()
	if err := smtpService.SendMail(
		[]string{input.To},
		"Forgot Password",
		htmlContent,
	); err != nil {
		return err
	}
	return nil
}

// NewMailService create new mail service and impl interface IMailService
func NewMailService() service.IMailService {
	return &MailService{}
}
