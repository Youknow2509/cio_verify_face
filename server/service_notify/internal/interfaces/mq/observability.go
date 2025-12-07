package mq

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// kafkaConsumerMetrics holds Prometheus metrics for Kafka consumption.
type kafkaConsumerMetrics struct {
	messagesTotal      *prometheus.CounterVec
	processingDuration *prometheus.HistogramVec
	errorsTotal        *prometheus.CounterVec
}

var (
	kafkaMetricsOnce sync.Once
	kafkaMetricsInst *kafkaConsumerMetrics
)

// getKafkaMetrics lazily registers and returns the Kafka consumer metrics.
func getKafkaMetrics() *kafkaConsumerMetrics {
	kafkaMetricsOnce.Do(func() {
		namespace := global.SettingServer.Server.Name
		kafkaMetricsInst = &kafkaConsumerMetrics{
			messagesTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: namespace,
					Name:      "kafka_messages_total",
					Help:      "Total Kafka messages consumed",
				},
				[]string{"topic", "status"},
			),
			processingDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace: namespace,
					Name:      "kafka_processing_duration_seconds",
					Help:      "Kafka message processing duration",
					Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
				},
				[]string{"topic", "status"},
			),
			errorsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: namespace,
					Name:      "kafka_errors_total",
					Help:      "Total Kafka consumer errors",
				},
				[]string{"topic", "phase"},
			),
		}

		prometheus.MustRegister(
			kafkaMetricsInst.messagesTotal,
			kafkaMetricsInst.processingDuration,
			kafkaMetricsInst.errorsTotal,
		)
	})
	return kafkaMetricsInst
}

// startKafkaSpan creates a span for Kafka message processing using the global tracer.
func startKafkaSpan(ctx context.Context, topic string, thread int) (context.Context, trace.Span) {
	tracer := otel.Tracer(global.SettingServer.Server.Name)
	return tracer.Start(ctx, "kafka.consume",
		trace.WithAttributes(
			attribute.String("kafka.topic", topic),
			attribute.Int("kafka.thread", thread),
		),
	)
}

// recordKafkaSuccess records successful processing metrics.
func recordKafkaSuccess(topic string, durationSeconds float64) {
	m := getKafkaMetrics()
	m.messagesTotal.WithLabelValues(topic, "success").Inc()
	m.processingDuration.WithLabelValues(topic, "success").Observe(durationSeconds)
}

// recordKafkaError records error metrics.
func recordKafkaError(topic, phase string, durationSeconds float64) {
	m := getKafkaMetrics()
	m.messagesTotal.WithLabelValues(topic, "error").Inc()
	m.processingDuration.WithLabelValues(topic, "error").Observe(durationSeconds)
	m.errorsTotal.WithLabelValues(topic, phase).Inc()
}
