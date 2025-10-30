package mail

import "errors"

/**
 * Interface for mail service for smtp
 */
type ISMTPService interface {
	SendMail(to []string, subject string, body string) error
}

/**
 * Manager instance for mail service
 */
var _vISMTPService ISMTPService

func SetSMTPService(svc ISMTPService) error {
	if svc == nil {
		return errors.New("smtp service is nil")
	}
	if _vISMTPService != nil {
		return errors.New("smtp service is already set")
	}
	_vISMTPService = svc
	return nil
}

func GetSMTPService() ISMTPService {
	return _vISMTPService
}
