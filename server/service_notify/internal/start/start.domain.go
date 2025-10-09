package start

import (
	domainMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mail"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mq"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/conn"
	infraMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/mail"
	infraMq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/mq"
)

func initDomain() error {
	// ============================================
	// 			Get client connection
	// ============================================
	postgres, err := infraConn.GetPostgresqlClient()
	if err != nil {
		return err
	}
	// ============================================
	// 			Initialize domain components
	// ============================================
	_ = postgres // TODO: use postgres connection here

	// ============================================
	// 			Initialize domain mail
	// ============================================
	implMailHtml := infraMail.NewHTMLContentMail()
	if err := domainMail.SetHtmlMailContent(implMailHtml); err != nil {
		return err
	}
	// ============================================
	smtpAuth, err := infraConn.GetSmtpClient()
	if err != nil {
		return err
	}
	implMailSmtp := infraMail.NewSMTPMail(
		smtpAuth,
		global.SettingServer.SMTP.Host,
		global.SettingServer.SMTP.Port,
	)
	if err := domainMail.SetSMTPService(implMailSmtp); err != nil {
		return err
	}
	// ============================================
	// 			Initialize domain kafka
	// ============================================
	implKafka := infraMq.NewKafkaReaderService(&global.SettingServer.Kafka)
	domainMq.InitKafkaReadService(implKafka)
	// v.v
	return nil
}
