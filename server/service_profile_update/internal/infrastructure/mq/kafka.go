package mq

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/mq"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/conn"
)

type KafkaWriterService struct {
	kafkaSetting *config.KafkaSetting
	kafkaTls     *tls.Config
	kafkaSasl    sasl.Mechanism
}

func (k *KafkaWriterService) WriteMessage(ctx context.Context, topic string, key string, value []byte) error {
	writer := k.getProducer()
	defer writer.Close()
	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

func (k *KafkaWriterService) WriteMessageRequireAck(ctx context.Context, topic string, key string, value []byte) error {
	writer := k.getProducerAckRequired()
	defer writer.Close()
	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

func (k *KafkaWriterService) WriteMessageRequireAllAck(ctx context.Context, topic string, key string, value []byte) error {
	writer := k.getProducerAllAckRequired()
	defer writer.Close()
	return writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

func (k *KafkaWriterService) getProducer() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: k.kafkaSetting.Brokers,
		Dialer: &kafka.Dialer{
			TLS:           k.kafkaTls,
			SASLMechanism: k.kafkaSasl,
		},
		BatchSize:    k.kafkaSetting.Producer.BatchSize,
		BatchBytes:   k.kafkaSetting.Producer.BatchBytes,
		ReadTimeout:  time.Duration(k.kafkaSetting.Producer.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(k.kafkaSetting.Producer.WriteTimeoutMs) * time.Millisecond,
		Async:        k.kafkaSetting.Producer.Async,
		Balancer:     getKafkaBalancer(k.kafkaSetting.Producer.Balancer),
		RequiredAcks: int(getKafkaRequiredAcks(constants.KAFKA_ACKS_NONE)),
	})
	writer.Compression = getKafkaCompression(k.kafkaSetting.Producer.CompressionType)
	return writer
}

func (k *KafkaWriterService) getProducerAckRequired() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: k.kafkaSetting.Brokers,
		Dialer: &kafka.Dialer{
			TLS:           k.kafkaTls,
			SASLMechanism: k.kafkaSasl,
		},
		BatchSize:    k.kafkaSetting.Producer.BatchSize,
		BatchBytes:   k.kafkaSetting.Producer.BatchBytes,
		ReadTimeout:  time.Duration(k.kafkaSetting.Producer.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(k.kafkaSetting.Producer.WriteTimeoutMs) * time.Millisecond,
		Async:        k.kafkaSetting.Producer.Async,
		Balancer:     getKafkaBalancer(k.kafkaSetting.Producer.Balancer),
		RequiredAcks: int(getKafkaRequiredAcks(constants.KAFKA_ACKS_LEADER)),
	})
	writer.Compression = getKafkaCompression(k.kafkaSetting.Producer.CompressionType)
	return writer
}

func (k *KafkaWriterService) getProducerAllAckRequired() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: k.kafkaSetting.Brokers,
		Dialer: &kafka.Dialer{
			TLS:           k.kafkaTls,
			SASLMechanism: k.kafkaSasl,
		},
		BatchSize:    k.kafkaSetting.Producer.BatchSize,
		BatchBytes:   k.kafkaSetting.Producer.BatchBytes,
		ReadTimeout:  time.Duration(k.kafkaSetting.Producer.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(k.kafkaSetting.Producer.WriteTimeoutMs) * time.Millisecond,
		Async:        k.kafkaSetting.Producer.Async,
		Balancer:     getKafkaBalancer(k.kafkaSetting.Producer.Balancer),
		RequiredAcks: int(getKafkaRequiredAcks(constants.KAFKA_ACKS_ALL)),
	})
	writer.Compression = getKafkaCompression(k.kafkaSetting.Producer.CompressionType)
	return writer
}

func getKafkaBalancer(balancerType int) kafka.Balancer {
	switch balancerType {
	case constants.KAFKA_BALANCER_ROUND_ROBIN:
		return &kafka.RoundRobin{}
	case constants.KAFKA_BALANCER_HASH:
		return &kafka.Hash{}
	default:
		return &kafka.RoundRobin{}
	}
}

func getKafkaRequiredAcks(acks int) kafka.RequiredAcks {
	switch acks {
	case constants.KAFKA_ACKS_NONE:
		return kafka.RequireNone
	case constants.KAFKA_ACKS_LEADER:
		return kafka.RequireOne
	case constants.KAFKA_ACKS_ALL:
		return kafka.RequireAll
	default:
		return kafka.RequireAll
	}
}

func getKafkaCompression(compression int) kafka.Compression {
	switch compression {
	case constants.KAFKA_COMPRESSION_GZIP:
		return kafka.Gzip
	case constants.KAFKA_COMPRESSION_SNAPPY:
		return kafka.Snappy
	case constants.KAFKA_COMPRESSION_LZ4:
		return kafka.Lz4
	case constants.KAFKA_COMPRESSION_ZSTD:
		return kafka.Zstd
	default:
		return kafka.Snappy
	}
}

func NewKafkaWriterService(kafkaSetting *config.KafkaSetting) domainMq.IKafkaWriter {
	conn.InitializeKafkaSecurity(kafkaSetting)
	kafkaTls, _ := conn.GetKafkaTls()
	kafkaSasl, _ := conn.GetKafkaSasl()
	return &KafkaWriterService{
		kafkaSetting: kafkaSetting,
		kafkaTls:     kafkaTls,
		kafkaSasl:    kafkaSasl,
	}
}
