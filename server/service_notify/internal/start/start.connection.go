package start

import (
	domainCache "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/cache"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/config"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/token"
	infraCache "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/cache"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/conn"
	infraToken "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/token"
)

var (
	_tokenService domainToken.ITokenService
)

func initConnectionToInfrastructure(setting *domainConfig.Setting) error {
	// initialize redis distributed cache
	if err := initRedisDistributedCache(&setting.Redis); err != nil {
		return err
	}
	// initialize posgresql
	if err := initConnectionPostgreSQL(&setting.Postgres); err != nil {
		return err
	}
	// initialize token service
	_tokenService = infraToken.NewTokenService(
		setting.JWT.Secret,
		setting.JWT.Issuer,
		setting.JWT.Subject,
		setting.JWT.Audience,
	)
	// initialize smtp
	if err := initSmtpClient(&setting.SMTP); err != nil {
		return err
	}
	// initialize kafka
	if err := initKafkaSecurity(&setting.Kafka); err != nil {
		return err
	}
	// v.v

	return nil
}

// Get TokenService returns the token service
func GetTokenService() domainToken.ITokenService {
	return _tokenService
}

func initSmtpClient(setting *domainConfig.SMTPSetting) error {
	err := infraConn.InitSmtpClient(setting)
	if err != nil {
		return err
	}
	return nil
}

func initKafkaSecurity(setting *domainConfig.KafkaSetting) error {
	err := infraConn.InitializeKafkaSecurity(setting)
	if err != nil {
		return err
	}
	return nil
}

func initRedisDistributedCache(setting *domainConfig.RedisSetting) error {
	distributedCacheImpl, err := infraCache.NewRedisDistributedCache(setting)
	if err != nil {
		return err
	}
	if err := domainCache.SetDistributedCache(distributedCacheImpl); err != nil {
		return err
	}
	return nil
}

func initConnectionPostgreSQL(setting *domainConfig.PostgresSetting) error {
	if err := infraConn.InitPostgresqlClient(setting); err != nil {
		return err
	}
	return nil
}
