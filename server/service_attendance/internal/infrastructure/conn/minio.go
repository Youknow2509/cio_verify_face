package clients

import (
	"errors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	constants "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/constants"
)

// minio client variables
var (
	vMinioClient *minio.Client
)

/**
 * Initializes the Minio client with the given settings.
 * @param settings *config.MinioSetting - The settings for the Minio client.
 * @return error - Returns an error if the initialization fails.
 */
func InitMinioClient(settings *domainConfig.MinioSetting) error {
	// Check if the Minio client is already initialized
	if vMinioClient != nil {
		return errors.New("minio client is already initialized")
	}
	var err error
	var client *minio.Client
	switch settings.Type {
	case constants.MinioTypeMinio:
		client, err = newMinioClient(settings)
	case constants.MinioTypeS3:
		settings.Endpoint = constants.S3Endpoint
		client, err = newMinioClient(settings)
	default:
		return errors.New("unsupported Minio client type")
	}
	if err != nil {
		return err
	}
	// Assign the client to the global variable
	vMinioClient = client
	return nil
}

/**
 * Gets the Minio client instance.
 * @return (*MinioClientSetting, error) - Returns the Minio client instance and an error if any.
 */
func GetMinioClient() (*minio.Client, error) {
	if vMinioClient == nil {
		return nil, errors.New("minio client is not initialized, please call InitMinioClient first")
	}
	return vMinioClient, nil
}

// =============================================================
//
//	Helper functions for Minio client
//
// =============================================================
// help create a new Minio client
func newMinioClient(settings *domainConfig.MinioSetting) (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(settings.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.AccessKey, settings.SecretKey, settings.Token),
		Secure: settings.EnableSSL,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}
