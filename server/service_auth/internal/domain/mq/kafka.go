package mq

import (
	"context"
	"errors"
)

// ====================================
// 			Kafka interface
// ====================================
type (
	// ================ Write message to Kafka ================
	IKafkaWrite interface {
		// Sync no requiretion
		WriteMessage(ctx context.Context, topic string, key string, value []byte) error
		// Sync requiretion a node kafka ack
		WriteMessageRequireAck(ctx context.Context, topic string, key string, value []byte) error
		// Sync requiretion all kafka ack
		WriteMessageRequireAllAck(ctx context.Context, topic string, key string, value []byte) error
	}

	// ================ Read message from Kafka ================
	IKafkaRead interface {
		// Sync read auto commit message after t time in group
		ReadMessageAutoCommit(ctx context.Context, topic string) (interface{}, error)
		// Sync read manual commit message in group
		ReadMessageManual(ctx context.Context, topic string, callback func(message interface{}) error) error
		// Sync read at offset message in group
		ReadMessageAtOffset(ctx context.Context, topic string, partition int32, offset int64) (interface{}, error)
		// Sync read from timestamp message in group
		ReadMessageFromTimestamp(ctx context.Context, topic string, partition int32, timestamp int64) (interface{}, error)

		// Sync read batch x messages manual commit in group
		ReadMessageBatchManual(ctx context.Context, topic string, partition int32, offset int64, limit int32, callback func(message interface{}) error) error
		// Sync read batch x messages at offset manual commit in group
		ReadMessageBatchAtOffsetManual(ctx context.Context, topic string, partition int32, offset int64, limit int32, callback func(message interface{}) error) error
		// Sync read batch x messages at timestamp manual commit in group
		ReadMessageBatchFromTimestampManual(ctx context.Context, topic string, partition int32, timestamp int64, limit int32, callback func(message interface{}) error) error

		// Listen topic 
		ReadListenTopicManual(ctx context.Context, topic string, callback func(message interface{}) error) error

		// ================ Help ================
		// Commit message
		CommitMessage(ctx context.Context, topic string, partition int32, offset int64) error
	}
)

/**
 * This variable is used to hold the instance of the Kafka service.
 */
var (
	vIKafkaWriterService IKafkaWrite
	vIKafkaReaderService IKafkaRead
)

/**
 * Initialize the Kafka Write service.
 */
func InitKafkaWriteService(service IKafkaWrite) {
	vIKafkaWriterService = service
}

/**
 * Initialize the Kafka Read service.
 */
func InitKafkaReadService(service IKafkaRead) {
	vIKafkaReaderService = service
}

/**
 * Get the Kafka service instance.
 * @return (instance, error)
 */
func GetKafkaWriteService() (IKafkaWrite, error) {
	if vIKafkaWriterService == nil {
		return nil, errors.New("kafka write service is not initialized, please call InitKafkaWriteService first")
	}
	return vIKafkaWriterService, nil
}

/**
 * Get the Kafka Read service instance.
 * @return (instance, error)
 */
func GetKafkaReadService() (IKafkaRead, error) {
	if vIKafkaReaderService == nil {
		return nil, errors.New("kafka read service is not initialized, please call InitKafkaReadService first")
	}
	return vIKafkaReaderService, nil
}

