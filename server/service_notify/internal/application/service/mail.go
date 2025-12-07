package service

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
)

/**
 * Mail service application
 */
type IMailService interface {
	SendMessageForgotPassword(
		ctx context.Context,
		input model.MailForgotPassword,
	) error
	SendPasswordResetNotification(
		ctx context.Context,
		input model.PasswordResetNotification,
	) error
	SendReportAttentionNotification(
		ctx context.Context,
		input model.ReportAttentionNotification,
	) error
}

/**
 * Manager instance of mail service
 */
var _vIMailService IMailService

func GetMailService() IMailService {
	return _vIMailService
}

func SetMailService(s IMailService) error {
	if s == nil {
		return errors.New("service init nil")
	}
	if _vIMailService != nil {
		return errors.New("service exists")
	}
	_vIMailService = s
	return nil
}