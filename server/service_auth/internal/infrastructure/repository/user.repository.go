package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/repository"
	db "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/gen"
)

/**
 * Struct impl IUserRepository
 */
type UserRepository struct {
	q db.Queries
}

// CreateUserSession implements repository.IUserRepository.
func (u *UserRepository) CreateUserSession(ctx context.Context, data *model.CreateUserSessionInput) error {
	return u.q.CreateUserSession(ctx, db.CreateUserSessionParams{
		SessionID:    pgtype.UUID{Bytes: data.SessionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: data.UserID, Valid: true},
		RefreshToken: data.RefreshToken,
		IpAddress:    &data.IPAddress,
		UserAgent:    pgtype.Text{String: data.UserAgent, Valid: true},
		ExpiresAt:    pgtype.Timestamptz{Time: data.ExpiredAt, Valid: true},
	})
}

// GetRefreshSessionInfo implements repository.IUserRepository.
func (u *UserRepository) GetRefreshSessionInfo(ctx context.Context, data *model.GetRefreshSessionInfoInput) (*model.GetRefreshSessionInfoOutput, error) {
	panic("unimplemented")
}

// GetUserBaseByEmail implements repository.IUserRepository.
func (u *UserRepository) GetUserBaseByEmail(ctx context.Context, email string) (*model.UserBaseInfoOutput, error) {
	response, err := u.q.GetUserBaseWithMail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &model.UserBaseInfoOutput{
		UserID:       response.UserID.String(),
		UserEmail:    response.Email,
		UserSalt:     response.Salt,
		UserPassword: response.PasswordHash,
		IsBlocked:    response.IsLocked.Bool,
		Role:         int(response.Role),
	}, nil
}

// GetUserBaseByID implements repository.IUserRepository.
func (u *UserRepository) GetUserBaseByID(ctx context.Context, userID uuid.UUID) (*model.UserBaseInfoOutput, error) {
	panic("unimplemented")
}

// RefreshSession implements repository.IUserRepository.
func (u *UserRepository) RefreshSession(ctx context.Context, data *model.RefreshSessionInput) error {
	panic("unimplemented")
}

// RemoveUserSession implements repository.IUserRepository.
func (u *UserRepository) RemoveUserSession(ctx context.Context, data *model.RemoveUserSessionInput) error {
	panic("unimplemented")
}

/**
 * New UserRepository
 */
func NewUserRepository(client *pgxpool.Pool) domainRepository.IUserRepository {
	return &UserRepository{
		q: *db.New(client),
	}
}
