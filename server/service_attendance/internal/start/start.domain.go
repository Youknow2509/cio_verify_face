package start

import (
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/token"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/repository"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/conn"
)

func initDomain() error {
	// ============================================
	// 			Get client connection
	// ============================================
	postgres, err := infraConn.GetPostgresqlClient()
	if err != nil {
		return err
	}
	cql, err := infraConn.GetScylladbClient()
	if err != nil {
		return err
	}
	// ============================================
	// 			Initialize domain components
	// ============================================
	// initialize ITokenService
	if err := domainToken.SetTokenService(
		GetTokenService(),
	); err != nil {
		return err
	}
	// init IAuditLogRepository
	if err := domainRepository.SetAuditRepository(
		infraRepository.NewAuditRepository(cql),
	); err != nil {
		return err
	}
	// init IAttendanceRepository
	if err := domainRepository.SetAttendanceRepository(
		infraRepository.NewAttendanceRepository(cql),
	); err != nil {
		return err
	}
	// init IUserRepository
	if err := domainRepository.SetUserRepository(
		infraRepository.NewUserRepository(postgres),
	); err != nil {
		return err
	}

	// v.v
	return nil
}
