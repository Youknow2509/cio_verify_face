package clients

import (
	"errors"
	"net/smtp"

	domainConfig "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/config"
)

// Smtp client variables
var (
	vSmtpAuth smtp.Auth
)

/**
 * Initializes the Smtp client.
 * @param *config.SmtpSetting - The configuration settings for Smtp.
 * @return error - Returns an error if the client initialization fails.
 */
func InitSmtpClient(smtpSetting *domainConfig.SMTPSetting) error {
	vSmtpAuth = smtp.PlainAuth(
		"",
		smtpSetting.Username,
		smtpSetting.Password,
		smtpSetting.Host,
	)
	return nil
}

/**
 * Get the Smtp client.
 * @return (smtp.Auth, error) - Returns the Smtp client and an error if it is not initialized.
 */
func GetSmtpClient() (smtp.Auth, error) {
	if vSmtpAuth == nil {
		return nil, errors.New("smtp client is not initialized, please call InitSmtpClient first")
	}
	return vSmtpAuth, nil
}

/**
 * NewSMTPClient
 */
func NewSMTPClient(smtpSetting *domainConfig.SMTPSetting) (smtp.Auth, error) {
	if err := InitSmtpClient(smtpSetting); err != nil {
		return nil, err
	}
	return GetSmtpClient()
}