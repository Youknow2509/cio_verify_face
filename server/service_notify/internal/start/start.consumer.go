package start

import (
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	interfacemq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/interfaces/mq"
)

// init consumer
func initConsumerKafka() error {
	// Existing notification topic listener
	mq := interfacemq.NewKafkaListenerData(
		constants.KAFKA_TOPIC_NOTIFICATION,
		global.SettingServer.Kafka.Consumer.Threads,
	)
	if err := mq.Listener(global.ContextSystem); err != nil {
		return err
	}

	// Password reset notifications listener - using separate config
	passwordResetListener := interfacemq.NewPasswordResetKafkaListener(
		global.SettingServer.PasswordResetNotifications.Topic,
		global.SettingServer.PasswordResetNotifications.Workers,
	)
	if err := passwordResetListener.Listener(global.ContextSystem); err != nil {
		return err
	}

	// Report attention notifications listener - using separate config
	reportAttentionListener := interfacemq.NewReportAttentionKafkaListener(
		global.SettingServer.ReportAttentionNotification.Topic,
		global.SettingServer.ReportAttentionNotification.Workers,
	)
	if err := reportAttentionListener.Listener(global.ContextSystem); err != nil {
		return err
	}

	return nil
}
