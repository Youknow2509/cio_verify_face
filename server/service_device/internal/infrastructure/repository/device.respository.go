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
 * Device repository implementation
 */
type DeviceRepository struct {
	db *database.Queries
}

// GetDeviceToken implements repository.IDeviceRepository.
func (d *DeviceRepository) GetDeviceToken(ctx context.Context, input *model.GetDeviceTokenInput) (*model.GetDeviceTokenOutput, error) {
	resp, err := d.db.GetDeviceToken(
		ctx,
		pgtype.UUID{Valid: true, Bytes: input.DeviceId},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.GetDeviceTokenOutput{
		DeviceId: input.DeviceId,
		Token:    resp,
	}, nil
}

// UpdateTokenDevice implements repository.IDeviceRepository.
func (d *DeviceRepository) UpdateTokenDevice(ctx context.Context, input *model.UpdateTokenDeviceInput) error {
	panic("unimplemented")
}

// DeviceExist implements repository.IDeviceRepository.
func (d *DeviceRepository) DeviceExist(ctx context.Context, input *model.DeviceExistInput) (bool, error) {
	_, err := d.db.CheckDeviceExist(
		ctx,
		pgtype.UUID{Valid: true, Bytes: input.DeviceId},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CreateNewDevice implements repository.IDeviceRepository.
func (d *DeviceRepository) CreateNewDevice(ctx context.Context, input *model.NewDevice) error {
	return d.db.CreateNewDevice(
		ctx,
		database.CreateNewDeviceParams{
			DeviceID:     pgtype.UUID{Valid: true, Bytes: input.DeviceId},
			CompanyID:    pgtype.UUID{Valid: true, Bytes: input.CompanyId},
			Name:         input.Name,
			Address:      pgtype.Text{Valid: true, String: input.Address},
			SerialNumber: pgtype.Text{Valid: true, String: input.SerialNumber},
			MacAddress:   pgtype.Text{Valid: true, String: input.MacAddress},
			Token:        input.Token,
		},
	)
}

// DeleteDevice implements repository.IDeviceRepository.
func (d *DeviceRepository) DeleteDevice(ctx context.Context, input *model.DeleteDeviceInput) error {
	return d.db.DeleteDevice(
		ctx,
		pgtype.UUID{Valid: true, Bytes: input.DeviceId},
	)
}

// DeviceInfo implements repository.IDeviceRepository.
func (d *DeviceRepository) DeviceInfo(ctx context.Context, input *model.DeviceInfoInput) (*model.DeviceInfoOutput, error) {
	resp, err := d.db.GetDeviceInfo(
		ctx,
		pgtype.UUID{Valid: true, Bytes: input.DeviceId},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.DeviceInfoOutput{
		DeviceId:        input.DeviceId,
		CompanyId:       resp.CompanyID.Bytes,
		Name:            resp.Name,
		Address:         resp.Address.String,
		SerialNumber:    resp.SerialNumber.String,
		MacAddress:      resp.MacAddress.String,
		IpAddress:       resp.IpAddress.String(),
		FirmwareVersion: resp.FirmwareVersion.String,
		LastHeartbeat:   resp.LastHeartbeat.Time.String(),
		Settings:        resp.Settings,
		CreateAt:        resp.CreatedAt.Time.String(),
		UpdateAt:        resp.UpdatedAt.Time.String(),
		Token:           resp.Token,
	}, nil
}

// DeviceInfoBase implements repository.IDeviceRepository.
func (d *DeviceRepository) DeviceInfoBase(ctx context.Context, input *model.DeviceInfoBaseInput) (*model.DeviceInfoBaseOutput, error) {
	resp, err := d.db.GetDeviceInfoBase(
		ctx,
		pgtype.UUID{Valid: true, Bytes: input.DeviceId},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.DeviceInfoBaseOutput{
		DeviceId:     input.DeviceId,
		CompanyId:    resp.CompanyID.Bytes,
		Name:         resp.Name,
		Address:      resp.Address.String,
		SerialNumber: resp.SerialNumber.String,
		MacAddress:   resp.MacAddress.String,
		Status:       int(resp.Status.Int16),
		CreateAt:     resp.CreatedAt.Time.String(),
		UpdateAt:     resp.UpdatedAt.Time.String(),
		Token:        resp.Token,
	}, nil
}

// DisableDevice implements repository.IDeviceRepository.
func (d *DeviceRepository) DisableDevice(ctx context.Context, input *model.DisableDeviceInput) error {
	return d.db.DisableDevice(ctx, pgtype.UUID{Valid: true, Bytes: input.DeviceId})
}

// EnableDevice implements repository.IDeviceRepository.
func (d *DeviceRepository) EnableDevice(ctx context.Context, input *model.EnableDeviceInput) error {
	return d.db.EnableDevice(ctx, pgtype.UUID{Valid: true, Bytes: input.DeviceId})
}

// ListDeviceInCompany implements repository.IDeviceRepository.
func (d *DeviceRepository) ListDeviceInCompany(ctx context.Context, input *model.ListDeviceInCompanyInput) (*model.ListDeviceInCompanyOutput, error) {
	resp, err := d.db.GetListDeviceInCompany(
		ctx,
		database.GetListDeviceInCompanyParams{
			CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyId},
			Limit:     int32(input.Limit),
			Offset:    int32(input.Offset),
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.ListDeviceInCompanyOutput{
				Devices: nil,
			}, nil
		}
	}
	devices := make([]*model.DeviceInfoBaseOutput, 0, len(resp))
	for _, device := range resp {
		devices = append(devices, &model.DeviceInfoBaseOutput{
			DeviceId:     device.DeviceID.Bytes,
			CompanyId:    device.CompanyID.Bytes,
			Name:         device.Name,
			Address:      device.Address.String,
			SerialNumber: device.SerialNumber.String,
			MacAddress:   device.MacAddress.String,
			Status:       int(device.Status.Int16),
			CreateAt:     device.CreatedAt.Time.String(),
			UpdateAt:     device.UpdatedAt.Time.String(),
			Token:        device.Token,
		})
	}
	return &model.ListDeviceInCompanyOutput{
		Devices: devices,
	}, nil
}

// UpdateDeviceInfo implements repository.IDeviceRepository.
func (d *DeviceRepository) UpdateDeviceInfo(ctx context.Context, input *model.UpdateDeviceInfoInput) error {
	return d.db.UpdateDeviceInfo(
		ctx,
		database.UpdateDeviceInfoParams{
			DeviceID:        pgtype.UUID{Valid: true, Bytes: input.DeviceId},
			FirmwareVersion: pgtype.Text{Valid: true, String: input.FirmwareVersion},
			SerialNumber:    pgtype.Text{Valid: true, String: input.SerialNumber},
			MacAddress:      pgtype.Text{Valid: true, String: input.MacAddress},
		},
	)
}

// UpdateDeviceLocation implements repository.IDeviceRepository.
func (d *DeviceRepository) UpdateDeviceLocation(ctx context.Context, input *model.UpdateDeviceLocationInput) error {
	return d.db.UpdateDeviceLocation(
		ctx,
		database.UpdateDeviceLocationParams{
			DeviceID:   pgtype.UUID{Valid: true, Bytes: input.DeviceId},
			LocationID: pgtype.UUID{Valid: true, Bytes: input.LocationId},
			Address:    pgtype.Text{Valid: true, String: input.Address},
		},
	)
}

// UpdateDeviceName implements repository.IDeviceRepository.
func (d *DeviceRepository) UpdateDeviceName(ctx context.Context, input *model.UpdateDeviceNameInput) error {
	return d.db.UpdateDeviceName(
		ctx,
		database.UpdateDeviceNameParams{
			DeviceID: pgtype.UUID{Valid: true, Bytes: input.DeviceId},
			Name:     input.Name,
		},
	)
}

// NewDeviceRepository create new instance and implement IDeviceRepository
func NewDeviceRepository(
	postgresConnect *pgxpool.Pool,
) domainRepo.IDeviceRepository {
	return &DeviceRepository{
		db: database.New(postgresConnect),
	}
}
