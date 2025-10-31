package start

import (
	domainCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/cache"
	libsConfig "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/config"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/mq"
	infraCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/cache"
	infraMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/mq"
)

func initConnectionToInfrastructure(setting *libsConfig.Setting) error {
	// initialize redis distributed cache
	if err := initRedisDistributedCache(&setting.Redis); err != nil {
		return err
	}
	// initialize kafka writer
	if err := initKafkaWriter(&setting.Kafka); err != nil {
		return err
	}
	// v.v

	return nil
}

func initKafkaWriter(setting *libsConfig.KafkaSetting) error {
	domainMq.InitKafkaWriteService(infraMq.NewKafkaWriterService(setting))
	return nil
}

func initRedisDistributedCache(setting *libsConfig.RedisSetting) error {
	distributedCacheImpl, err := infraCache.NewRedisDistributedCache(setting)
	if err != nil {
		return err
	}
	if err := domainCache.SetDistributedCache(distributedCacheImpl); err != nil {
		return err
	}
	return nil
}
