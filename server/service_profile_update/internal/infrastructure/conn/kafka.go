package conn

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
)

var (
	vKafkaTls  *tls.Config
	vKafkaSasl sasl.Mechanism
)

func InitializeKafkaSecurity(settings *config.KafkaSetting) error {
	if settings == nil {
		return errors.New("kafka settings are nil")
	}
	if err := setKafkaTls(&settings.TLS); err != nil {
		return fmt.Errorf("failed to set Kafka TLS configuration: %w", err)
	}
	if err := setKafkaSasl(&settings.SASL); err != nil {
		return fmt.Errorf("failed to set Kafka SASL configuration: %w", err)
	}
	return nil
}

func GetKafkaTls() (*tls.Config, error) {
	if vKafkaTls == nil {
		return nil, nil // TLS is optional
	}
	return vKafkaTls, nil
}

func GetKafkaSasl() (sasl.Mechanism, error) {
	if vKafkaSasl == nil {
		return nil, nil // SASL is optional
	}
	return vKafkaSasl, nil
}

func setKafkaTls(setting *config.KafkaTLSSetting) error {
	if setting == nil || !setting.Enabled {
		vKafkaTls = nil
		return nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: setting.SkipVerify,
	}

	if setting.CAFile != "" {
		caCert, err := os.ReadFile(setting.CAFile)
		if err != nil {
			return fmt.Errorf("failed to read CA file: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return fmt.Errorf("failed to append CA certs from %s", setting.CAFile)
		}
		tlsConfig.RootCAs = caCertPool
	}

	vKafkaTls = tlsConfig
	return nil
}

func setKafkaSasl(setting *config.KafkaSASLSetting) error {
	if setting == nil || !setting.Enabled {
		vKafkaSasl = nil
		return nil
	}

	switch setting.Mechanism {
	case constants.KAFKA_SASL_MECHANISM_PLAIN:
		vKafkaSasl = plain.Mechanism{
			Username: setting.Username,
			Password: setting.Password,
		}
	case constants.KAFKA_SASL_MECHANISM_SCRAM_SHA256:
		mech, err := scram.Mechanism(scram.SHA256, setting.Username, setting.Password)
		if err != nil {
			return fmt.Errorf("failed to create SCRAM-SHA256 mechanism: %w", err)
		}
		vKafkaSasl = mech
	case constants.KAFKA_SASL_MECHANISM_SCRAM_SHA512:
		mech, err := scram.Mechanism(scram.SHA512, setting.Username, setting.Password)
		if err != nil {
			return fmt.Errorf("failed to create SCRAM-SHA512 mechanism: %w", err)
		}
		vKafkaSasl = mech
	default:
		vKafkaSasl = nil
	}

	return nil
}
