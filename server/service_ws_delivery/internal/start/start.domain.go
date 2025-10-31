package start

import (
	libsClients "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/conn"
	domainHealth "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/health"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/mq"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/repository"
	infraHealth "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/health"
	infraMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/mq"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/repository"
)

func initDomain() error {
	// init manager connetion session
	redisClient, err := libsClients.GetRedisClient()
	if err != nil {
		return err
	}
	implRedisManagerConnectionRepository := infraRepository.NewRedisManagerConnectionRepository(
		redisClient,
	)
	if err := domainRepository.SetManagerConnectionRepository(implRedisManagerConnectionRepository); err != nil {
		return err
	}
	// init send event client to kafka
	implSendEventClient := infraMq.NewSendEventToKafka()
	if err := domainMq.SetSendEventToKafka(implSendEventClient); err != nil {
		return err
	}
	// init healthy check
	implHealthCheck := infraHealth.NewHealthCheck()
	if err := domainHealth.SetHealthCheck(implHealthCheck); err != nil {
		return err
	}
	// v.v
	return nil
}
