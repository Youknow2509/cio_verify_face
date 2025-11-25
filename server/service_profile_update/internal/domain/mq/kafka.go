package mq

import (
	"context"
	"errors"
)

// =================================
// Kafka Writer Interface:
// =================================
type IKafkaWriter interface {
	// Write message without requiring ack
	WriteMessage(ctx context.Context, topic string, key string, value []byte) error

	// Write message requiring leader ack
	WriteMessageRequireAck(ctx context.Context, topic string, key string, value []byte) error

	// Write message requiring all replicas ack
	WriteMessageRequireAllAck(ctx context.Context, topic string, key string, value []byte) error
}

// =================================
// Kafka Reader Interface:
// =================================
type IKafkaReader interface {
	// Read message with auto commit
	ReadMessageAutoCommit(ctx context.Context, topic string) (interface{}, error)

	// Read message with manual commit
	ReadMessageManual(ctx context.Context, topic string, callback func(message interface{}) error) error

	// Listen to topic
	ReadListenTopicManual(ctx context.Context, topic string, callback func(message interface{}) error) error

	// Commit message
	CommitMessage(ctx context.Context, topic string, partition int32, offset int64) error
}

// =================================
// Kafka Variables:
// =================================
var (
	_kafkaWriter IKafkaWriter
	_kafkaReader IKafkaReader
)

// =================================
// Setters and Getters:
// =================================
func SetKafkaWriter(writer IKafkaWriter) error {
	if _kafkaWriter != nil {
		return errors.New("kafka writer already initialized")
	}
	_kafkaWriter = writer
	return nil
}

func GetKafkaWriter() (IKafkaWriter, error) {
	if _kafkaWriter == nil {
		return nil, errors.New("kafka writer not initialized")
	}
	return _kafkaWriter, nil
}

func SetKafkaReader(reader IKafkaReader) error {
	if _kafkaReader != nil {
		return errors.New("kafka reader already initialized")
	}
	_kafkaReader = reader
	return nil
}

func GetKafkaReader() (IKafkaReader, error) {
	if _kafkaReader == nil {
		return nil, errors.New("kafka reader not initialized")
	}
	return _kafkaReader, nil
}
