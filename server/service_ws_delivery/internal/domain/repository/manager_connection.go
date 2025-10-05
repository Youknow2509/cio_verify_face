package repository

import (
	"context"
	"errors"

	model "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
)

// ======================================================================================================
//
//	Cache Repository Interface - handle cache client connection
//
// ======================================================================================================
type IManagerConnectionRepository interface {
	CreateConnection(ctx context.Context, input *model.CreateConnectionInput) (bool, error)
	RemoveConnection(ctx context.Context, input *model.RemoveConnectionInput) (bool, error)
}

// variable to hold the repository implementation
var (
	_vIManagerConnectionRepository IManagerConnectionRepository
)

// ======================================================================================================
//
//	Getter and setter for the repository implementation
//
// ======================================================================================================
func GetManagerConnectionRepository() IManagerConnectionRepository {
	return _vIManagerConnectionRepository
}

func SetManagerConnectionRepository(managerConnectionRepository IManagerConnectionRepository) error {
	if managerConnectionRepository == nil {
		return errors.New("manager connection repository cannot be nil")
	}
	if _vIManagerConnectionRepository != nil {
		return errors.New("manager connection repository is already set")
	}
	_vIManagerConnectionRepository = managerConnectionRepository
	return nil
}
