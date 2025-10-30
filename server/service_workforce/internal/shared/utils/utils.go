package utils

import (
	"github.com/segmentio/kafka-go"
	constants "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/constants"
)

// =================================
//
//	Kafka utils
//
// =================================
func GetKafkaBalancer(balancerType int) kafka.Balancer {
	// 	Use Case										Balancer nên dùng
	//  Cần đảm bảo thứ tự theo key (user, order)		Hash hoặc Murmur2
	//  Hệ thống nhiều ngôn ngữ cần giống Kafka Java	Murmur2Balancer hoặc ReferenceHash
	//  Chia đều message, không cần theo key			RoundRobin
	//  Giảm tải động giữa partition					LeastBytes
	//  ....
	switch balancerType {
	case constants.KAFKA_BALANCER_ROUND_ROBIN:
		return &kafka.RoundRobin{}
	case constants.KAFKA_BALANCER_LEAST_BYTES:
		return &kafka.LeastBytes{}
	case constants.KAFKA_BALANCER_HASH:
		return &kafka.Hash{}
	case constants.KAFKA_BALANCER_REFERENCE_HASH:
		return &kafka.ReferenceHash{}
	case constants.KAFKA_BALANCER_CRC32:
		return &kafka.CRC32Balancer{}
	case constants.KAFKA_BALANCER_MURMUR2:
		return &kafka.Murmur2Balancer{}
	case constants.KAFKA_BALANCER_CUSTOM:
		panic("Custom balancer is not implemented yet")
	default:
		return &kafka.RoundRobin{}
	}
}

func GetKafkaRequiredAcks(acks int) kafka.RequiredAcks {
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

func GetKafkaCompression(compression int) kafka.Compression {
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
		return kafka.Snappy // Use Snappy as a default, or choose another supported compression
	}
}

func GetKafkaOffset(offsetType int) int64 {
	switch offsetType {
	case constants.KAFKA_AUTO_OFFSET_RESET_EARLIEST:
		return kafka.FirstOffset
	case constants.KAFKA_AUTO_OFFSET_RESET_LATEST:
		return kafka.LastOffset
	default:
		return kafka.LastOffset
	}
}
