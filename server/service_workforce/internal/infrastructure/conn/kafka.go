package clients

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/config"
)

// Kafka tls, sasl variable
var (
	vKafkaTls  *tls.Config
	vKafkaSasl sasl.Mechanism
)

/**
 * Initialize Kafka TLS and SASL configurations.
 * @param settings - Kafka settings
 * @return *errors.Error - returns an error if the initialization fails
 */
func InitializeKafkaSecurity(settings *config.KafkaSetting) error {
	if settings == nil {
		return errors.New("kafka settings are nil")
	}
	// set tls configuration
	if err := setKafkaTls(&settings.TLS); err != nil {
		return errors.New("failed to set Kafka TLS configuration")
	}
	// set sasl configuration
	if err := setKafkaSasl(&settings.SASL); err != nil {
		return errors.New("failed to set Kafka SASL configuration")
	}
	return nil
}

/**
 * Get Kafka TLS security configuration.
 * @return *tls.Config - returns the TLS configuration for Kafka
 */
func GetKafkaTls() (*tls.Config, error) {
	if vKafkaTls == nil {
		return nil, errors.New("kafka TLS configuration is not initialized, please call InitializeKafkaSecurity first")
	}
	return vKafkaTls, nil
}

/**
 * Get Kafka SASL security configuration.
 * @return sasl.Mechanism - returns the SASL mechanism for Kafka
 */
func GetKafkaSasl() (sasl.Mechanism, error) {
	if vKafkaSasl == nil {
		return nil, errors.New("kafka SASL configuration is not initialized, please call InitializeKafkaSecurity first")
	}
	return vKafkaSasl, nil
}

// ===========================================================
//
//	Kafka helpers
//
// ===========================================================
// help set security TLS – Transport Layer Security for kafka
func setKafkaTls(setting *config.KafkaTLSSetting) error {
	if setting == nil || !setting.Enabled {
		vKafkaTls = nil
		return nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: setting.SkipVerify,
	}
	// nếu có CA file, đọc và add vào RootCAs
	if setting.CAFile != "" {
		caCert, err := os.ReadFile(setting.CAFile)
		if err != nil {
			return err
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

// help set security SASL – Simple Authentication and Security Layer for kafka
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
			return err
		}
		vKafkaSasl = mech
	case constants.KAFKA_SASL_MECHANISM_SCRAM_SHA512:
		mech, err := scram.Mechanism(scram.SHA512, setting.Username, setting.Password)
		if err != nil {
			return err
		}
		vKafkaSasl = mech
	default:
		vKafkaSasl = nil
	}
	return nil
}
