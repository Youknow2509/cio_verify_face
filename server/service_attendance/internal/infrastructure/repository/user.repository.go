package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	db "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/gen"
)

// ============================================
// User repository impl
// ============================================
type UserRepository struct {
	q db.Queries
}

// UserIsManagerCompany implements repository.IUserRepository.
func (u *UserRepository) UserIsManagerCompany(ctx context.Context, input *domainModel.UserIsManagerCompanyInput) (bool, error) {
	_, err := u.q.CheckUserIsManagementInCompany(
		ctx,
		db.CheckUserIsManagementInCompanyParams{
			CompanyID:  pgtype.UUID{Valid: true, Bytes: input.CompanyID},
			EmployeeID: pgtype.UUID{Valid: true, Bytes: input.UserID},
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetCompanyIdUser implements repository.IUserRepository.
func (u *UserRepository) GetCompanyIdUser(ctx context.Context, input *domainModel.GetCompanyIdUserInput) (*domainModel.GetCompanyIdUserOutput, error) {
	reps, err := u.q.GetCompanyIdUser(
		ctx,
		pgtype.UUID{Valid: true, Bytes: input.UserID},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &domainModel.GetCompanyIdUserOutput{
		CompanyID: reps.Bytes,
	}, nil
}

// UserIsEmployeeInCompany implements repository.IUserRepository.
func (u *UserRepository) UserIsEmployeeInCompany(ctx context.Context, input *domainModel.UserIsEmployeeInCompanyInput) (bool, error) {
	_, err := u.q.GetCompanyIdUser(
		ctx,
		pgtype.UUID{Valid: true, Bytes: input.UserID},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// New instance user repository and impl IUserRepository
func NewUserRepository(conn *pgxpool.Pool) domainRepo.IUserRepository {
	return &UserRepository{
		q: *db.New(conn),
	}
}
