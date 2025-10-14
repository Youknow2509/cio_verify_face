package start

import (
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/token"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/conn"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/repository"
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
	// initialize ITokenService
	if err := domainToken.SetTokenService(
		GetTokenService(),
	); err != nil {
		return err
	}
	// init ICompanyRepository
	if err := domainRepository.SetCompanyRepository(
		infraRepository.NewCompanyRepository(postgres),
	); err != nil {
		return err
	}
	// init IAuditLogRepository
	if err := domainRepository.SetAuditRepository(
		infraRepository.NewAuditRepository(postgres),
	); err != nil {
		return err
	}

	// v.v
	return nil
}
