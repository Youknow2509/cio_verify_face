package impl

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	constants "github.com/youknow2509/cio_verify_face/server/service_notify/internal/constants"
	config "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/config"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mq"
	clients "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/conn"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/shared/utils"
)

type (
	// ===== Kafka Writer Service =====
	KafkaWriterService struct {
		kafkaSetting *config.KafkaSetting
		kafkaTls     *tls.Config
		kafkaSasl    sasl.Mechanism
	}

	// ===== Kafka Reader Service =====
	KafkaReaderService struct {
		kafkaSetting *config.KafkaSetting
		kafkaTls     *tls.Config
		kafkaSasl    sasl.Mechanism
	}
)

// ==================== KafkaReaderService methods ====================

// ReadListenTopicManual implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadListenTopicManual(ctx context.Context, topic string, callback func(message []byte) error) error {
	reader := k.getConsumer(topic)
	defer reader.Close()
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := callback(m.Value); err != nil {
			return err
		}
		// Commit the message after processing
		if err := reader.CommitMessages(ctx, m); err != nil {
			return err
		}
	}
}

// ReadMessageAtOffset implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadMessageAtOffset(ctx context.Context, topic string, partition int32, offset int64) ([]byte, error) {
	reader := k.getConsumer(topic)
	defer reader.Close()
	reader.SetOffset(offset)
	m, err := reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return m.Value, nil
}

// ReadMessageAutoCommit implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadMessageAutoCommit(ctx context.Context, topic string) ([]byte, error) {
	reader := k.getConsumerAutoCommit(topic)
	m, err := reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return m.Value, nil
}

// ReadMessageBatchAtOffsetManual implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadMessageBatchAtOffsetManual(ctx context.Context, topic string, partition int32, offset int64, limit int32, callback func(message []byte) error) error {
	reader := k.getConsumer(topic)
	defer reader.Close()
	reader.SetOffset(offset)
	for i := int32(0); i < limit; i++ {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := callback(m.Value); err != nil {
			return err
		}
		// Commit the message after processing
		if err := reader.CommitMessages(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

// ReadMessageBatchFromTimestampManual implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadMessageBatchFromTimestampManual(ctx context.Context, topic string, partition int32, timestamp int64, limit int32, callback func(message []byte) error) error {
	reader := k.getConsumer(topic)
	defer reader.Close()
	seekTime := time.UnixMilli(timestamp)
	reader.SetOffsetAt(ctx, seekTime)
	for i := int32(0); i < limit; i++ {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := callback(m.Value); err != nil {
			return err
		}
	}
	return nil
}

// ReadMessageBatchManual implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadMessageBatchManual(ctx context.Context, topic string, partition int32, offset int64, limit int32, callback func(message []byte) error) error {
	return k.ReadMessageBatchAtOffsetManual(ctx, topic, partition, offset, limit, callback)
}

// ReadMessageFromTimestamp implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadMessageFromTimestamp(ctx context.Context, topic string, partition int32, timestamp int64) ([]byte, error) {
	reader := k.getConsumer(topic)
	defer reader.Close()
	seekTime := time.UnixMilli(timestamp)
	reader.SetOffsetAt(ctx, seekTime)
	m, err := reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return m.Value, nil
}

// ReadMessageManual implements domainMq.IKafkaRead.
func (k *KafkaReaderService) ReadMessageManual(ctx context.Context, topic string, callback func(message []byte) error) error {
	reader := k.getConsumer(topic)
	defer reader.Close()
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := callback(m.Value); err != nil {
			return err
		}
	}
}

// commitMessage implements domainMq.IKafkaRead.
func (k *KafkaReaderService) CommitMessage(ctx context.Context, topic string, partition int32, offset int64) error {
	cl := k.getConsumer(topic)
	defer cl.Close()
	message := &kafka.Message{
		Topic:     topic,
		Partition: int(partition),
		Offset:    offset,
	}
	return cl.CommitMessages(ctx, *message)
}

// ==================== KafkaWriterService methods ====================

// WriteMessage implements domainMq.IKafkaWrite.
func (k *KafkaWriterService) WriteMessage(ctx context.Context, topic string, key string, value []byte) error {
	writer := k.getProducer()
	defer writer.Close()
	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

// WriteMessageRequireAck implements domainMq.IKafkaWrite.
func (k *KafkaWriterService) WriteMessageRequireAck(ctx context.Context, topic string, key string, value []byte) error {
	writer := k.getProducerAckRequired()
	defer writer.Close()
	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

// WriteMessageRequireAllAck implements domainMq.IKafkaWrite.
func (k *KafkaWriterService) WriteMessageRequireAllAck(ctx context.Context, topic string, key string, value []byte) error {
	writer := k.getProducerAllAckRequired()
	defer writer.Close()
	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

// =============================================================
//
//	NewKafkaService creates a new KafkaService instance
//
// =============================================================
func NewKafkaWriterService(kafkaSetting *config.KafkaSetting) domainMq.IKafkaWrite {
	clients.InitializeKafkaSecurity(kafkaSetting)
	kafkaTls, _ := clients.GetKafkaTls()
	kafkaSasl, _ := clients.GetKafkaSasl()
	return &KafkaWriterService{
		kafkaSetting: kafkaSetting,
		kafkaTls:     kafkaTls,
		kafkaSasl:    kafkaSasl,
	}
}

// =============================================================

// NewKafkaReaderService creates a new KafkaReaderService instance
func NewKafkaReaderService(kafkaSetting *config.KafkaSetting) domainMq.IKafkaRead {
	clients.InitializeKafkaSecurity(kafkaSetting)
	kafkaTls, _ := clients.GetKafkaTls()
	kafkaSasl, _ := clients.GetKafkaSasl()
	return &KafkaReaderService{
		kafkaSetting: kafkaSetting,
		kafkaTls:     kafkaTls,
		kafkaSasl:    kafkaSasl,
	}
}

// ===== Helper Functions =====
func (k *KafkaWriterService) getProducer() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: k.kafkaSetting.Brokers,
		Dialer: &kafka.Dialer{
			TLS:           k.kafkaTls,
			SASLMechanism: k.kafkaSasl,
		},
		// Producer configuration
		BatchSize:    k.kafkaSetting.Producer.BatchSize,
		BatchBytes:   k.kafkaSetting.Producer.BatchBytes,
		ReadTimeout:  time.Duration(k.kafkaSetting.Producer.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(k.kafkaSetting.Producer.WriteTimeoutMs) * time.Millisecond,
		Async:        k.kafkaSetting.Producer.Async,
		// Balancer configuration
		Balancer: utils.GetKafkaBalancer(k.kafkaSetting.Producer.Balancer),
		// Compression configuration
		CompressionCodec: nil,
		// Required acks configuration
		RequiredAcks: int(utils.GetKafkaRequiredAcks(constants.KAFKA_ACKS_NONE)),
	})
	writer.Compression = utils.GetKafkaCompression(k.kafkaSetting.Producer.CompressionType)
	return writer
}

func (k *KafkaWriterService) getProducerAckRequired() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: k.kafkaSetting.Brokers,
		Dialer: &kafka.Dialer{
			TLS:           k.kafkaTls,
			SASLMechanism: k.kafkaSasl,
		},
		// Producer configuration
		BatchSize:    k.kafkaSetting.Producer.BatchSize,
		BatchBytes:   k.kafkaSetting.Producer.BatchBytes,
		ReadTimeout:  time.Duration(k.kafkaSetting.Producer.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(k.kafkaSetting.Producer.WriteTimeoutMs) * time.Millisecond,
		Async:        k.kafkaSetting.Producer.Async,
		// Balancer configuration
		Balancer: utils.GetKafkaBalancer(k.kafkaSetting.Producer.Balancer),
		// Compression configuration
		CompressionCodec: nil,
		// Required acks configuration
		RequiredAcks: int(utils.GetKafkaRequiredAcks(constants.KAFKA_ACKS_LEADER)),
	})
	writer.Compression = utils.GetKafkaCompression(k.kafkaSetting.Producer.CompressionType)
	return writer
}

func (k *KafkaWriterService) getProducerAllAckRequired() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: k.kafkaSetting.Brokers,
		Dialer: &kafka.Dialer{
			TLS:           k.kafkaTls,
			SASLMechanism: k.kafkaSasl,
		},
		// Producer configuration
		BatchSize:    k.kafkaSetting.Producer.BatchSize,
		BatchBytes:   k.kafkaSetting.Producer.BatchBytes,
		ReadTimeout:  time.Duration(k.kafkaSetting.Producer.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(k.kafkaSetting.Producer.WriteTimeoutMs) * time.Millisecond,
		Async:        k.kafkaSetting.Producer.Async,
		// Balancer configuration
		Balancer: utils.GetKafkaBalancer(k.kafkaSetting.Producer.Balancer),
		// Compression configuration
		CompressionCodec: nil,
		// Required acks configuration
		RequiredAcks: int(utils.GetKafkaRequiredAcks(constants.KAFKA_ACKS_ALL)),
	})
	writer.Compression = utils.GetKafkaCompression(k.kafkaSetting.Producer.CompressionType)
	return writer
}

func (k *KafkaReaderService) getConsumer(topic string) *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: k.kafkaSetting.Brokers,
			GroupID: k.kafkaSetting.Consumer.GroupID,
			Topic:   topic,
			Dialer: &kafka.Dialer{
				TLS:           k.kafkaTls,
				SASLMechanism: k.kafkaSasl,
			},
			CommitInterval:    0, // Disable auto-commit
			MinBytes:          k.kafkaSetting.Consumer.MinBytes,
			MaxBytes:          k.kafkaSetting.Consumer.MaxBytes,
			MaxWait:           time.Duration(k.kafkaSetting.Consumer.MaxWaitMs) * time.Millisecond,
			ReadBatchTimeout:  time.Duration(k.kafkaSetting.Consumer.ReadBatchTimeoutMs) * time.Millisecond,
			HeartbeatInterval: time.Duration(k.kafkaSetting.Consumer.HeartbeatIntervalMs) * time.Millisecond,
			SessionTimeout:    time.Duration(k.kafkaSetting.Consumer.SessionTimeoutMs) * time.Millisecond,
			RebalanceTimeout:  time.Duration(k.kafkaSetting.Consumer.RebalanceTimeoutMs) * time.Millisecond,
			JoinGroupBackoff:  time.Duration(k.kafkaSetting.Consumer.JoinGroupBackoffMs) * time.Millisecond,
			ReadLagInterval:   time.Duration(k.kafkaSetting.Consumer.ReadLagIntervalMs) * time.Millisecond,
			MaxAttempts:       k.kafkaSetting.Consumer.MaxAttempts,
			QueueCapacity:     k.kafkaSetting.Consumer.QueueCapacity,
			RetentionTime:     time.Duration(k.kafkaSetting.Consumer.RetentionTimeMs) * time.Millisecond,
		})
}

func (k *KafkaReaderService) getConsumerAutoCommit(topic string) *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: k.kafkaSetting.Brokers,
			Topic:   topic,
			GroupID: k.kafkaSetting.Consumer.GroupID,
			Dialer: &kafka.Dialer{
				TLS:           k.kafkaTls,
				SASLMechanism: k.kafkaSasl,
			},
			CommitInterval:    time.Duration(k.kafkaSetting.Consumer.CommitIntervalMs) * time.Millisecond, // Enable auto-commit
			MinBytes:          k.kafkaSetting.Consumer.MinBytes,
			MaxBytes:          k.kafkaSetting.Consumer.MaxBytes,
			MaxWait:           time.Duration(k.kafkaSetting.Consumer.MaxWaitMs) * time.Millisecond,
			ReadBatchTimeout:  time.Duration(k.kafkaSetting.Consumer.ReadBatchTimeoutMs) * time.Millisecond,
			HeartbeatInterval: time.Duration(k.kafkaSetting.Consumer.HeartbeatIntervalMs) * time.Millisecond,
			SessionTimeout:    time.Duration(k.kafkaSetting.Consumer.SessionTimeoutMs) * time.Millisecond,
			RebalanceTimeout:  time.Duration(k.kafkaSetting.Consumer.RebalanceTimeoutMs) * time.Millisecond,
			JoinGroupBackoff:  time.Duration(k.kafkaSetting.Consumer.JoinGroupBackoffMs) * time.Millisecond,
			ReadLagInterval:   time.Duration(k.kafkaSetting.Consumer.ReadLagIntervalMs) * time.Millisecond,
			MaxAttempts:       k.kafkaSetting.Consumer.MaxAttempts,
			QueueCapacity:     k.kafkaSetting.Consumer.QueueCapacity,
			RetentionTime:     time.Duration(k.kafkaSetting.Consumer.RetentionTimeMs) * time.Millisecond,
		})
}
