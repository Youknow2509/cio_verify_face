package start

import (
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/token"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/repository"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/conn"
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
	// init IDeviceRepository
	if err := domainRepository.SetDeviceRepository(
		infraRepository.NewDeviceRepository(postgres),
	); err != nil {
		return err
	}
	// v.v
	return nil
}
