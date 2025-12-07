package start

import (
	domainFace "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/face"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/token"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/conn"
	infraFace "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/face"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/repository"
	infraToken "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/token"
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
	// initialize IDeviceRepository
	if err := domainRepository.SetDeviceRepository(
		infraRepository.NewDeviceRepository(postgres),
	); err != nil {
		return err
	}
	// initialize token service
	if err := domainToken.SetTokenService(
		infraToken.NewTokenService(authGrpcClient),
	); err != nil {
		return err
	}
	// initialize face verification service
	if err := domainFace.SetFaceVerificationService(
		infraFace.NewFaceVerificationService(faceGrpcClient),
	); err != nil {
		return err
	}
	// ============================================

	// v.v
	return nil
}
