package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/repository"
	database "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/infrastructure/gen"
)

/**
 * Shift repository implementation
 */
type ShiftRepository struct {
	db   *database.Queries
	pool *pgxpool.Pool
}

// DisableShiftWithId implements repository.IShiftRepository.
func (s *ShiftRepository) DisableShiftWithId(ctx context.Context, input *model.DisableShiftInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}
	err := s.db.DisableShiftWithId(ctx, database.DisableShiftWithIdParams{
		ShiftID:   pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyId},
	})
	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New("no shift found with the given ID and company ID")
	default:
		return err
	}
}

// EnableShiftWithId implements repository.IShiftRepository.
func (s *ShiftRepository) EnableShiftWithId(ctx context.Context, input *model.EnableShiftInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}
	err := s.db.EnableShiftWithId(ctx, database.EnableShiftWithIdParams{
		ShiftID:   pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyId},
	})
	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New("no shift found with the given ID and company ID")
	default:
		return err
	}
}

// NewShiftRepository create new instance and implement IShiftRepository
func NewShiftRepository(
	postgresConnect *pgxpool.Pool,
) domainRepo.IShiftRepository {
	return &ShiftRepository{
		db:   database.New(postgresConnect),
		pool: postgresConnect,
	}
}

// Helper functions for type conversion
func toPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func toPgTime(t time.Time) pgtype.Time {
	// pgtype.Time stores microseconds since midnight
	hour, min, sec := t.Clock()
	microseconds := int64(hour*3600+min*60+sec) * 1000000
	microseconds += int64(t.Nanosecond() / 1000)
	return pgtype.Time{Microseconds: microseconds, Valid: true}
}

func toPgInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func fromPgText(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

func fromPgTime(t pgtype.Time) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	// Convert microseconds since midnight to time.Time
	microseconds := t.Microseconds
	hours := microseconds / 3600000000
	minutes := (microseconds % 3600000000) / 60000000
	seconds := (microseconds % 60000000) / 1000000
	nanoseconds := (microseconds % 1000000) * 1000

	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(),
		int(hours), int(minutes), int(seconds), int(nanoseconds), now.Location())
}

func fromPgInt4(i pgtype.Int4) int32 {
	if i.Valid {
		return i.Int32
	}
	return 0
}

func fromPgBool(b pgtype.Bool) bool {
	if b.Valid {
		return b.Bool
	}
	return false
}

func fromPgTimestamptz(t pgtype.Timestamptz) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

// CreateShift implements repository.IShiftRepository.
func (s *ShiftRepository) CreateShift(ctx context.Context, input *model.CreateShiftInput) (uuid.UUID, error) {
	if input == nil {
		return uuid.UUID{}, errors.New("input cannot be nil")
	}

	id, err := s.db.CreateShift(ctx, database.CreateShiftParams{
		CompanyID:             pgtype.UUID{Valid: true, Bytes: input.CompanyID},
		Name:                  input.Name,
		Description:           toPgText(input.Description),
		StartTime:             toPgTime(input.StartTime),
		EndTime:               toPgTime(input.EndTime),
		BreakDurationMinutes:  toPgInt4(input.BreakDurationMinutes),
		GracePeriodMinutes:    toPgInt4(input.GracePeriodMinutes),
		EarlyDepartureMinutes: toPgInt4(input.EarlyDepartureMinutes),
		WorkDays:              input.WorkDays,
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	if !id.Valid {
		return uuid.UUID{}, errors.New("failed to create shift: empty id returned")
	}
	return id.Bytes, nil
}

// ListShifts implements repository.IShiftRepository.
func (s *ShiftRepository) ListShifts(ctx context.Context, input *model.ListShiftsInput) ([]*model.Shift, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	rows, err := s.db.ListShifts(ctx, database.ListShiftsParams{
		CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyID},
		Column2:   input.IsActive,
		Limit:     input.Limit,
		Offset:    input.Offset,
	})
	if err != nil {
		return nil, err
	}

	out := make([]*model.Shift, 0, len(rows))
	for _, r := range rows {
		out = append(out, &model.Shift{
			ShiftID:               r.ShiftID.Bytes,
			CompanyID:             r.CompanyID.Bytes,
			Name:                  r.Name,
			Description:           fromPgText(r.Description),
			StartTime:             fromPgTime(r.StartTime),
			EndTime:               fromPgTime(r.EndTime),
			BreakDurationMinutes:  fromPgInt4(r.BreakDurationMinutes),
			GracePeriodMinutes:    fromPgInt4(r.GracePeriodMinutes),
			EarlyDepartureMinutes: fromPgInt4(r.EarlyDepartureMinutes),
			WorkDays:              r.WorkDays,
			IsFlexible:            fromPgBool(r.IsFlexible),
			OvertimeAfterMinutes:  fromPgInt4(r.OvertimeAfterMinutes),
			IsActive:              fromPgBool(r.IsActive),
			CreatedAt:             fromPgTimestamptz(r.CreatedAt),
			UpdatedAt:             fromPgTimestamptz(r.UpdatedAt),
		})
	}
	return out, nil
}

// GetShiftByID implements repository.IShiftRepository.
func (s *ShiftRepository) GetShiftByID(ctx context.Context, shiftID uuid.UUID) (*model.Shift, error) {
	r, err := s.db.GetShiftByID(ctx, pgtype.UUID{Valid: true, Bytes: shiftID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &model.Shift{
		ShiftID:               r.ShiftID.Bytes,
		CompanyID:             r.CompanyID.Bytes,
		Name:                  r.Name,
		Description:           fromPgText(r.Description),
		StartTime:             fromPgTime(r.StartTime),
		EndTime:               fromPgTime(r.EndTime),
		BreakDurationMinutes:  fromPgInt4(r.BreakDurationMinutes),
		GracePeriodMinutes:    fromPgInt4(r.GracePeriodMinutes),
		EarlyDepartureMinutes: fromPgInt4(r.EarlyDepartureMinutes),
		WorkDays:              r.WorkDays,
		IsFlexible:            fromPgBool(r.IsFlexible),
		OvertimeAfterMinutes:  fromPgInt4(r.OvertimeAfterMinutes),
		IsActive:              fromPgBool(r.IsActive),
		CreatedAt:             fromPgTimestamptz(r.CreatedAt),
		UpdatedAt:             fromPgTimestamptz(r.UpdatedAt),
	}, nil
}

// UpdateTimeShift implements repository.IShiftRepository.
func (s *ShiftRepository) UpdateTimeShift(ctx context.Context, input *model.UpdateTimeShiftInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	return s.db.UpdateTimeShift(ctx, database.UpdateTimeShiftParams{
		ShiftID:               pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		StartTime:             toPgTime(input.StartTime),
		EndTime:               toPgTime(input.EndTime),
		BreakDurationMinutes:  toPgInt4(input.BreakDurationMinutes),
		GracePeriodMinutes:    toPgInt4(input.GracePeriodMinutes),
		EarlyDepartureMinutes: toPgInt4(input.EarlyDepartureMinutes),
		WorkDays:              input.WorkDays,
	})
}

// DeleteShift implements repository.IShiftRepository.
func (s *ShiftRepository) DeleteShift(ctx context.Context, shiftID uuid.UUID) error {
	if shiftID == uuid.Nil {
		return errors.New("shiftID cannot be empty")
	}
	return s.db.DeleteShift(ctx, pgtype.UUID{Valid: true, Bytes: shiftID})
}

// GetShiftsIdForCompany implements repository.IShiftRepository.
func (s *ShiftRepository) GetShiftsIdForCompany(ctx context.Context, companyID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	if companyID == uuid.Nil {
		return nil, errors.New("companyID cannot be empty")
	}

	rows, err := s.db.GetShiftsIdForCompany(ctx, database.GetShiftsIdForCompanyParams{
		CompanyID: pgtype.UUID{Valid: true, Bytes: companyID},
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, id := range rows {
		if id.Valid {
			ids = append(ids, id.Bytes)
		}
	}
	return ids, nil
}
