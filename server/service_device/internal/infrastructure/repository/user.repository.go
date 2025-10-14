package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/repository"
	db "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/gen"
)

/**
 * Struct impl IUserRepository
 */
type UserRepository struct {
	q db.Queries
}

// GetUserInfoByID implements repository.IUserRepository.
func (u *UserRepository) GetUserInfoByID(ctx context.Context, userID uuid.UUID) (*model.UserInfoOutput, error) {
	response, err := u.q.GetUserInfoWithID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &model.UserInfoOutput{
		Email:     response.Email,
		Phone:     response.Phone,
		FullName:  response.FullName,
		AvatarURL: response.AvatarUrl.String,
	}, nil
}

// GetUserSessionByID implements repository.IUserRepository.
func (u *UserRepository) GetUserSessionByID(ctx context.Context, sessionID uuid.UUID) (*model.UserSessionOutput, error) {
	response, err := u.q.GetUserSessionByID(
		ctx,
		pgtype.UUID{Bytes: sessionID, Valid: true},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &model.UserSessionOutput{
		SessionID:    response.SessionID.Bytes,
		UserID:       response.UserID.Bytes,
		RefreshToken: response.RefreshToken,
		IPAddress:    *response.IpAddress,
		UserAgent:    response.UserAgent.String,
		CreatedAt:    response.CreatedAt.Time,
		ExpiredAt:    response.ExpiresAt.Time,
	}, nil
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
	return u.q.UpdateUserSession(
		ctx,
		db.UpdateUserSessionParams{
			SessionID:    pgtype.UUID{Bytes: data.SessionID, Valid: true},
			RefreshToken: data.RefreshToken,
			ExpiresAt:    pgtype.Timestamptz{Time: data.ExpiredAt, Valid: true},
		},
	)
}

// RemoveUserSession implements repository.IUserRepository.
func (u *UserRepository) RemoveUserSession(ctx context.Context, data *model.RemoveUserSessionInput) error {
	return u.q.DeleteUserSessionByID(ctx, pgtype.UUID{Bytes: data.SessionID, Valid: true})
}

/**
 * New UserRepository
 */
func NewUserRepository(client *pgxpool.Pool) domainRepository.IUserRepository {
	return &UserRepository{
		q: *db.New(client),
	}
}
