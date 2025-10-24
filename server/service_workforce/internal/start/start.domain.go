package start

import (
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/token"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/infrastructure/conn"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/infrastructure/repository"
	infraToken "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/infrastructure/token"
)

func initDomain() error {
	// ============================================
	// 			Get client connection
	// ============================================
	postgres, err := infraConn.GetPostgresqlClient()
	if err != nil {
		return err
	}
	// ============================================
	// 			Initialize domain components
	// ============================================
	// initialize IUserRepository
	if err := domainRepository.SetUserRepository(
		infraRepository.NewUserRepository(postgres),
	); err != nil {
		return err
	}
	// initialize token service
	if err := domainToken.SetTokenService(
		infraToken.NewTokenService(grpcClient),
	); err != nil {
		return err
	}
	// ============================================

	// v.v
	return nil
}
