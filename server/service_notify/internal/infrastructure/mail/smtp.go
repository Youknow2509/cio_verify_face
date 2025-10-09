package mail

import (
	"fmt"
	"net/smtp"

	domainMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mail"
)

/**
 * Struct smtp service
 */
type SMTPMail struct {
	auth smtp.Auth
	from string
	host string
	port int
}

// SendMail implements mail.ISMTPService.
func (s *SMTPMail) SendMail(to []string, subject string, body string) error {
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n", to, subject, body,
	))
	addr := s.host + ":" + string(rune(s.port))
	if err := smtp.SendMail(
		addr,
		s.auth,
		s.from,
		to,
		msg,
	); err != nil {
		return err
	}
	return nil
}

/**
 * New smtp mail and impl interface ISMTPService
 */
func NewSMTPMail(
	auth smtp.Auth, host string, port int,
) domainMail.ISMTPService {
	return &SMTPMail{
		auth: auth,
		host: host,
		port: port,
	}
}
