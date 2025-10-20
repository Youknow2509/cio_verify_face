package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/repository"
	database "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/gen"
)

/**
 * User repository implementation
 */
type UserRepository struct {
	db *database.Queries
}

// UserExistsInCompany implements repository.IUserRepository.
func (u *UserRepository) UserExistsInCompany(ctx context.Context, input *model.UserExistsInCompanyInput) (bool, error) {
	exists, err := u.db.CheckUserExistInCompany(
		ctx,
		database.CheckUserExistInCompanyParams{
			EmployeeID: pgtype.UUID{Valid: true, Bytes: input.UserID},
			CompanyID:  pgtype.UUID{Valid: true, Bytes: input.CompanyID},
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

// NewUserRepository create new instance and implement IUserRepository
func NewUserRepository(
	postgresConnect *pgxpool.Pool,
) domainRepo.IUserRepository {
	return &UserRepository{
		db: database.New(postgresConnect),
	}
}
