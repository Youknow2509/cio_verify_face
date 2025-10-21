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

// Device repository implementations
type DeviceRepository struct {
	q db.Queries
}

// CheckTokenDevice implements repository.IDeviceRepository.
func (d *DeviceRepository) CheckTokenDevice(ctx context.Context, input *model.CheckTokenDeviceInput) (bool, string, error) {
	resp, err := d.q.CheckTokenDevice(
		ctx,
		db.CheckTokenDeviceParams{
			DeviceID: pgtype.UUID{Valid: true, Bytes: input.DeviceId},
			Token:    input.Token,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "", nil
		}
		return false, "", err
	}
	if int(resp.Int16) == 1 {
		return true, "", nil
	}
	return false, "", nil
}

// DeviceExist implements repository.IDeviceRepository.
func (d *DeviceRepository) DeviceExist(ctx context.Context, deviceID uuid.UUID) (bool, error) {
	exist, err := d.q.DeviceExists(
		ctx,
		pgtype.UUID{
			Valid: true,
			Bytes: deviceID,
		},
	)
	if err != nil {
		return false, err
	}
	if exist.Bytes == deviceID {
		return true, nil
	}
	return false, nil
}

// BlockDeviceToken implements repository.IDeviceRepository.
func (d *DeviceRepository) BlockDeviceToken(ctx context.Context, deviceToken uuid.UUID) error {
	return d.q.BlockDeviceToken(
		ctx,
		pgtype.UUID{
			Valid: true,
			Bytes: deviceToken,
		},
	)
}

// CreateDeviceToken implements repository.IDeviceRepository.
func (d *DeviceRepository) CreateDeviceToken(ctx context.Context, input *model.CreateDeviceTokenInput) error {
	err := d.q.CreateDeviceToken(
		ctx,
		db.CreateDeviceTokenParams{
			DeviceID: pgtype.UUID{
				Valid: true,
				Bytes: input.DeviceId,
			},
			Token: input.NewToken,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

/**
 * New DeviceRepository
 */
func NewDeviceRepository(client *pgxpool.Pool) domainRepository.IDeviceRepository {
	return &DeviceRepository{
		q: *db.New(client),
	}
}
