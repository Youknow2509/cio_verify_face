package impl

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service"
	domainMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mail"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
)

// Structure mail service
type MailService struct {
}

// SendMessageForgotPassword implements service.IMailService.
func (m *MailService) SendMessageForgotPassword(ctx context.Context, input model.MailForgotPassword) error {
	if global.Logger != nil {
		global.Logger.Info("sending forgot password email", "to", input.To)
	}
	htmlGen := domainMail.GetHtmlMailContent()
	htmlContent, err := htmlGen.ForgotPassword(
		input.To,
		input.UrlAuth,
		input.NewPassword,
		input.Expired,
	)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("generate forgot password email failed", "error", err)
		}
		return err
	}
	smtpService := domainMail.GetSMTPService()
	if err := smtpService.SendMail(
		[]string{input.To},
		"Forgot Password",
		htmlContent,
	); err != nil {
		if global.Logger != nil {
			global.Logger.Error("send forgot password email failed", "error", err)
		}
		return err
	}
	if global.Logger != nil {
		global.Logger.Info("forgot password email sent", "to", input.To)
	}
	return nil
}

// SendPasswordResetNotification implements service.IMailService.
func (m *MailService) SendPasswordResetNotification(ctx context.Context, input model.PasswordResetNotification) error {
	if global.Logger != nil {
		global.Logger.Info("sending password reset notification", "to", input.To)
	}
	htmlGen := domainMail.GetHtmlMailContent()
	htmlContent, err := htmlGen.PasswordResetNotification(
		input.To,
		input.FullName,
		input.ResetURL,
		input.ExpiresIn,
	)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("generate password reset notification failed", "error", err)
		}
		return err
	}
	smtpService := domainMail.GetSMTPService()
	if err := smtpService.SendMail(
		[]string{input.To},
		"Password Reset Request",
		htmlContent,
	); err != nil {
		if global.Logger != nil {
			global.Logger.Error("send password reset notification failed", "error", err)
		}
		return err
	}
	if global.Logger != nil {
		global.Logger.Info("password reset notification sent", "to", input.To)
	}
	return nil
}

// SendReportAttentionNotification implements service.IMailService.
func (m *MailService) SendReportAttentionNotification(ctx context.Context, input model.ReportAttentionNotification) error {
	if global.Logger != nil {
		global.Logger.Info("sending report attention notification", "email", input.Email)
	}
	htmlGen := domainMail.GetHtmlMailContent()
	htmlContent, err := htmlGen.ReportAttentionNotification(
		input.Email,
		input.DownloadURL,
		input.Type,
		input.Format,
		input.StartDate,
		input.EndDate,
	)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("generate report attention notification failed", "error", err)
		}
		return err
	}
	smtpService := domainMail.GetSMTPService()
	if err := smtpService.SendMail(
		[]string{input.Email},
		"Your Report is Ready",
		htmlContent,
	); err != nil {
		if global.Logger != nil {
			global.Logger.Error("send report attention notification failed", "error", err)
		}
		return err
	}
	if global.Logger != nil {
		global.Logger.Info("report attention notification sent", "email", input.Email)
	}
	return nil
}

// NewMailService create new mail service and impl interface IMailService
func NewMailService() service.IMailService {
	return &MailService{}
}
