package repository

import (
	"context"
	"errors"
	model "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
)

/**
 * Interface for audit repository
 */
type IAuditRepository interface {
	AddAuditLog(ctx context.Context, log *model.AuditLog) error
}

/**
 * Variable for Audit repository instance
 */
var _vAuditRepository IAuditRepository

/**
 * Set the Audit repository instance
 */
func SetAuditRepository(v IAuditRepository) error {
	if _vAuditRepository != nil {
		return errors.New("Audit repository initialization failed, not nil")
	}
	_vAuditRepository = v
	return nil
}

/**
 * Get the Audit repository instance
 */
func GetAuditRepository() (IAuditRepository, error) {
	if _vAuditRepository == nil {
		return nil, errors.New("Audit repository not initialized")
	}
	return _vAuditRepository, nil
}
