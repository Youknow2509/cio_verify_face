package repository

import (
	"context"
	"errors"

	model "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
)

// ============================================
// User repository interface
// ============================================
type IUserRepository interface {
	GetListTimeShiftEmployee(ctx context.Context, input *model.GetListTimeShiftEmployeeInput) ([]model.ShiftTimeEmployee, error)
	UserIsManagerCompany(ctx context.Context, input *model.UserIsManagerCompanyInput) (bool, error)
	UserIsEmployeeInCompany(ctx context.Context, input *model.UserIsEmployeeInCompanyInput) (bool, error)
	GetCompanyIdUser(ctx context.Context, input *model.GetCompanyIdUserInput) (*model.GetCompanyIdUserOutput, error)
}

// ============================================
// Manager instance user repository
// ============================================
var _vIUserRepository IUserRepository

// Getter instance user repository
func GetUserRepository() IUserRepository {
	return _vIUserRepository
}

// Setter instance user repository
func SetUserRepository(repo IUserRepository) error {
	if repo == nil {
		return errors.New("user repository set is nil")
	}
	if _vIUserRepository != nil {
		return errors.New("user repository already set")
	}
	_vIUserRepository = repo
	return nil
}
