package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/repository"
	db "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/gen"
)

/**
 * Struct impl ICompanyRepository
 */
type CompanyRepository struct {
	q db.Queries
}

// GetCompanyUser implements repository.ICompanyRepository.
func (c *CompanyRepository) GetCompanyUser(ctx context.Context, input *model.GetCompanyUserInput) (*model.GetCompanyUserOutput, error) {
	response, err := c.q.GetCompanyUser(
		ctx,
		pgtype.UUID{
			Valid: true,
			Bytes: input.UserID,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &model.GetCompanyUserOutput{
		CompanyID: response.Bytes,
	}, nil
}

// CheckUserIsManagementInCompany implements repository.ICompanyRepository.
func (c *CompanyRepository) CheckUserIsManagementInCompany(ctx context.Context, data *model.CheckCompanyIsManagementInCompanyInput) (bool, error) {
	response, err := c.q.CheckUserIsManagementInCompany(
		ctx,
		db.CheckUserIsManagementInCompanyParams{
			CompanyID: pgtype.UUID{
				Valid: true,
				Bytes: data.CompanyID,
			},
			EmployeeID: pgtype.UUID{
				Valid: true,
				Bytes: data.UserID,
			},
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if response == "" {
		return false, nil
	}
	return true, nil
}

// CheckDeviceExistsInCompany implements repository.ICompanyRepository.
func (c *CompanyRepository) CheckDeviceExistsInCompany(ctx context.Context, data *model.CheckDeviceExistsInCompanyInput) (bool, error) {
	response, err := c.q.CheckDeviceExistInCompany(
		ctx,
		db.CheckDeviceExistInCompanyParams{
			CompanyID: pgtype.UUID{
				Valid: true,
				Bytes: data.CompanyID,
			},
			DeviceID: pgtype.UUID{
				Valid: true,
				Bytes: data.DeviceID,
			},
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if response == "" {
		return false, nil
	}
	return true, nil
}

// DeleteDeviceSession implements repository.ICompanyRepository.
func (c *CompanyRepository) DeleteDeviceSession(ctx context.Context, data *model.DeleteDeviceSessionInput) error {
	return c.q.DeleteDeviceSessionByID(
		ctx,
		pgtype.UUID{
			Valid: true,
			Bytes: data.DeviceId,
		},
	)
}

// UpdateDeviceSession implements repository.ICompanyRepository.
func (c *CompanyRepository) UpdateDeviceSession(ctx context.Context, data *model.UpdateDeviceSessionInput) error {
	return c.q.UpdateDeviceSession(
		ctx,
		db.UpdateDeviceSessionParams{
			DeviceID: pgtype.UUID{
				Valid: true,
				Bytes: data.DeviceId,
			},
			Token: data.Token,
		},
	)
}

/**
 * New CompanyRepository
 */
func NewCompanyRepository(client *pgxpool.Pool) domainRepository.ICompanyRepository {
	return &CompanyRepository{
		q: *db.New(client),
	}
}
