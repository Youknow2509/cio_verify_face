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
 * ShiftUser repository implementation
 */
type ShiftUserRepository struct {
	db   *database.Queries
	pool *pgxpool.Pool
}

// AddListShiftForEmployees implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) AddListShiftForEmployees(ctx context.Context, input *model.AddListShiftForEmployeesInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	listError := make([]error, 0)
	for _, employeeId := range input.EmployeeIDs {
		err := s.db.AddShiftForEmployee(ctx, database.AddShiftForEmployeeParams{
			EmployeeID:    pgtype.UUID{Valid: true, Bytes: employeeId},
			ShiftID:       pgtype.UUID{Valid: true, Bytes: input.ShiftID},
			EffectiveFrom: toPgDate(input.EffectiveFrom),
			EffectiveTo:   toPgDate(input.EffectiveTo),
		})
		if err != nil {
			listError = append(listError, err)
		}
	}
	if len(listError) > 0 {
		return errors.New("one or more errors occurred while adding shifts for employees")
	}
	return nil
}

// NewShiftUserRepository create new instance and implement IShiftUserRepository
func NewShiftUserRepository(
	postgresConnect *pgxpool.Pool,
) domainRepo.IShiftUserRepository {
	return &ShiftUserRepository{
		db:   database.New(postgresConnect),
		pool: postgresConnect,
	}
}

// GetShiftEmployeeWithEffectiveDate implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) GetShiftEmployeeWithEffectiveDate(ctx context.Context, input *model.GetShiftEmployeeWithEffectiveDateInput) ([]*model.EmployeeShiftRow, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	rows, err := s.db.GetShiftEmployeeWithEffectiveDate(ctx, database.GetShiftEmployeeWithEffectiveDateParams{
		EmployeeID:    pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
		EffectiveFrom: toPgDate(input.EffectiveFrom),
		Limit:         input.Limit,
		Offset:        input.Offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*model.EmployeeShiftRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, &model.EmployeeShiftRow{
			EmployeeShiftID: r.EmployeeShiftID.Bytes,
			ShiftID:         r.ShiftID.Bytes,
			EffectiveFrom:   fromPgDate(r.EffectiveFrom),
			EffectiveTo:     fromPgDate(r.EffectiveTo),
			IsActive:        fromPgBoolValue(r.IsActive),
		})
	}
	return out, nil
}

// EditEffectiveShiftForEmployee implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) EditEffectiveShiftForEmployee(ctx context.Context, input *model.EditEffectiveShiftForEmployeeInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	return s.db.EditEffectiveShiftForEmployee(ctx, database.EditEffectiveShiftForEmployeeParams{
		EmployeeShiftID: pgtype.UUID{Valid: true, Bytes: input.EmployeeShiftID},
		EffectiveFrom:   toPgDate(input.EffectiveFrom),
		EffectiveTo:     toPgDate(input.EffectiveTo),
	})
}

// DeleteEmployeeShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) DeleteEmployeeShift(ctx context.Context, employeeShiftID uuid.UUID) error {
	if employeeShiftID == uuid.Nil {
		return errors.New("employeeShiftID cannot be empty")
	}
	return s.db.DeleteEmployeeShift(ctx, pgtype.UUID{Valid: true, Bytes: employeeShiftID})
}

// DisableEmployeeShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) DisableEmployeeShift(ctx context.Context, employeeShiftID uuid.UUID) error {
	if employeeShiftID == uuid.Nil {
		return errors.New("employeeShiftID cannot be empty")
	}
	return s.db.DisableEmployeeShift(ctx, pgtype.UUID{Valid: true, Bytes: employeeShiftID})
}

// EnableEmployeeShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) EnableEmployeeShift(ctx context.Context, employeeShiftID uuid.UUID) error {
	if employeeShiftID == uuid.Nil {
		return errors.New("employeeShiftID cannot be empty")
	}
	return s.db.EnableEmployeeShift(ctx, pgtype.UUID{Valid: true, Bytes: employeeShiftID})
}

// AddShiftForEmployee implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) AddShiftForEmployee(ctx context.Context, input *model.AddShiftForEmployeeInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}
	if input.EmployeeID == uuid.Nil {
		return errors.New("employeeID cannot be empty")
	}
	if input.ShiftID == uuid.Nil {
		return errors.New("shiftID cannot be empty")
	}

	return s.db.AddShiftForEmployee(ctx, database.AddShiftForEmployeeParams{
		EmployeeID:    pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
		ShiftID:       pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		EffectiveFrom: toPgDate(input.EffectiveFrom),
		EffectiveTo:   toPgDate(input.EffectiveTo),
	})
}

// CheckUserExistShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) CheckUserExistShift(ctx context.Context, input *model.CheckUserExistShiftInput) (bool, error) {
	if input == nil {
		return false, errors.New("input cannot be nil")
	}

	_, err := s.db.CheckUserExistShift(ctx, database.CheckUserExistShiftParams{
		EmployeeID:    pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
		EffectiveFrom: toPgDate(input.EffectiveFrom),
		EffectiveTo:   toPgDate(input.EffectiveTo),
		Limit:         input.Limit,
		Offset:        input.Offset,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Helper functions for type conversion
func toPgTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func fromPgTimestamp(t pgtype.Timestamptz) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

func toPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func fromPgDate(d pgtype.Date) time.Time {
	if d.Valid {
		return d.Time
	}
	return time.Time{}
}

func toPgBoolValue(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func fromPgBoolValue(b pgtype.Bool) bool {
	if b.Valid {
		return b.Bool
	}
	return false
}
