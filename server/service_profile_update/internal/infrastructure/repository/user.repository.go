package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/repository"
	db "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/gen"
)

/**
 * Struct impl IUserRepository
 */
type UserRepository struct {
	q db.Queries
}

// GetUserByID implements repository.IUserRepository
func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	result, err := r.q.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.UserInfo{
		UserID:    result.UserID.Bytes,
		Email:     result.Email,
		Phone:     result.Phone,
		FullName:  result.FullName,
		AvatarURL: result.AvatarUrl.String,
		Role:      model.Role(result.Role),
	}, nil
}

// GetUserByEmail implements repository.IUserRepository
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.UserInfo, error) {
	result, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.UserInfo{
		UserID:    result.UserID.Bytes,
		Email:     result.Email,
		Phone:     result.Phone,
		FullName:  result.FullName,
		AvatarURL: result.AvatarUrl.String,
		Role:      model.Role(result.Role),
	}, nil
}

// UserBelongsToCompany implements repository.IUserRepository
func (r *UserRepository) UserBelongsToCompany(ctx context.Context, userID, companyID uuid.UUID) (bool, error) {
	result, err := r.q.CheckUserBelongsToCompany(ctx, db.CheckUserBelongsToCompanyParams{
		EmployeeID: pgtype.UUID{Bytes: userID, Valid: true},
		CompanyID:  pgtype.UUID{Bytes: companyID, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

// GetEmployeeInfo implements repository.IUserRepository
func (r *UserRepository) GetEmployeeInfo(ctx context.Context, employeeID uuid.UUID) (*model.EmployeeInfo, error) {
	result, err := r.q.GetEmployeeInfo(ctx, pgtype.UUID{Bytes: employeeID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.EmployeeInfo{
		EmployeeID:   result.EmployeeID.Bytes,
		CompanyID:    result.CompanyID.Bytes,
		EmployeeCode: result.EmployeeCode,
		Department:   result.Department.String,
		Position:     result.Position.String,
	}, nil
}

// UpdateUserPassword implements repository.IUserRepository
func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, salt, passwordHash string) error {
	return r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		Salt:         salt,
		PasswordHash: passwordHash,
	})
}

// IsCompanyAdmin implements repository.IUserRepository
func (r *UserRepository) IsCompanyAdmin(ctx context.Context, userID, companyID uuid.UUID) (bool, error) {
	result, err := r.q.IsCompanyAdmin(ctx, db.IsCompanyAdminParams{
		EmployeeID: pgtype.UUID{Bytes: userID, Valid: true},
		CompanyID:  pgtype.UUID{Bytes: companyID, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

/**
 * NewUserRepository creates a new repository
 */
func NewUserRepository(client *pgxpool.Pool) domainRepository.IUserRepository {
	return &UserRepository{
		q: *db.New(client),
	}
}
