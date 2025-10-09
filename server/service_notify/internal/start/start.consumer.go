package start

import (
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	interfacemq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/interfaces/mq"
)

// init consumer
func initConsumerKafka() error {
	mq := interfacemq.NewKafkaListenerData(
		constants.KAFKA_TOPIC_NOTIFICATION,
		global.SettingServer.Kafka.Consumer.Threads,
	)
	if err := mq.Listener(global.ContextSystem); err != nil {
		return err
	}
	return nil
}
