package repository

import (
	"context"
	"errors"
	"fmt"
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

// GetListShiftForEmployee returns paginated shift assignments for a single employee.
func (s *ShiftUserRepository) GetListShiftForEmployee(ctx context.Context, input *model.GetListShiftForEmployeeInput) (*model.GetListShiftForEmployeeOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	// Build base query with optional company filter to ensure tenant isolation.
	baseQuery := "FROM employee_shifts es JOIN work_shifts ws ON ws.shift_id = es.shift_id WHERE es.employee_id = $1"
	args := []interface{}{input.EmployeeID}
	nextArg := 2
	if input.CompanyID != uuid.Nil {
		baseQuery = fmt.Sprintf("%s AND ws.company_id = $%d", baseQuery, nextArg)
		args = append(args, input.CompanyID)
		nextArg++
	}

	// Count total rows for pagination.
	countQuery := "SELECT COUNT(1) " + baseQuery
	var total int64
	if err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	// Fetch current page.
	listQuery := fmt.Sprintf(`SELECT es.shift_id, ws.name, ws.start_time::text, ws.end_time::text, es.effective_from, es.effective_to, es.is_active %s ORDER BY es.effective_from DESC LIMIT $%d OFFSET $%d`, baseQuery, nextArg, nextArg+1)
	argsWithPaging := append(args, input.Limit, input.Offset)
	rows, err := s.pool.Query(ctx, listQuery, argsWithPaging...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shifts := make([]*model.ShiftInfoForEmployee, 0)
	for rows.Next() {
		var (
			shiftID  pgtype.UUID
			name     string
			startStr string
			endStr   string
			effFrom  pgtype.Date
			effTo    pgtype.Date
			isActive pgtype.Bool
		)
		if scanErr := rows.Scan(&shiftID, &name, &startStr, &endStr, &effFrom, &effTo, &isActive); scanErr != nil {
			return nil, scanErr
		}
		shifts = append(shifts, &model.ShiftInfoForEmployee{
			ShiftID:        shiftID.Bytes,
			ShiftName:      name,
			ShiftStartTime: startStr,
			ShiftEndTime:   endStr,
			EffectiveFrom:  fromPgDate(effFrom),
			EffectiveTo:    fromPgDate(effTo),
			IsActive:       fromPgBoolValue(isActive),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &model.GetListShiftForEmployeeOutput{
		Shifts:   shifts,
		Total:    int32(total),
		PageSize: input.Limit,
	}, nil
}

const queryGetShiftEmployeeAll = `
SELECT 
	shift_id,
	effective_from,
	effective_to,
	is_active
FROM employee_shifts
WHERE employee_id = $1
ORDER BY effective_from DESC
LIMIT $2 OFFSET $3`

// DeleteListEmployeeShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) DeleteListEmployeeShift(ctx context.Context, input *model.DeleteListEmployeeShiftInput) (string, error) {
	if input == nil {
		return "", errors.New("input cannot be nil")
	}
	errStr := "Failed to delete shift for employee ID s:\n"
	for _, employeeId := range input.EmployeeIDs {
		err := s.db.DeleteEmployeeShift(ctx, database.DeleteEmployeeShiftParams{
			EmployeeID: pgtype.UUID{Valid: true, Bytes: employeeId},
			ShiftID:    pgtype.UUID{Valid: true, Bytes: input.ShiftId},
		})
		if err != nil {
			errStr += "- " + employeeId.String() + "\n"
		}
	}
	return errStr, nil
}

// GetListEmployeeDonotInShift implements repository.IShiftUserRepository using sqlc.
func (s *ShiftUserRepository) GetListEmployeeDonotInShift(ctx context.Context, input *model.GetListEmployyeShiftInput) (*model.GetListEmployyeShiftOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}
	// Count total
	total, err := s.db.CountEmployeesDonotInShiftCurrent(ctx, database.CountEmployeesDonotInShiftCurrentParams{
		CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyID},
		ShiftID:   pgtype.UUID{Valid: true, Bytes: input.ShiftID},
	})
	if err != nil {
		return nil, err
	}
	// Fetch page
	rows, err := s.db.GetListEmployeeDonotInShift(ctx, database.GetListEmployeeDonotInShiftParams{
		CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyID},
		ShiftID:   pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		Limit:     input.Limit,
		Offset:    input.Offset,
	})
	if err != nil {
		return nil, err
	}
	out := &model.GetListEmployyeShiftOutput{
		EmployeeIDs: make([]*model.EmployeeShiftInfoBase, 0, len(rows)),
		Total:       int32(total),
		PageSize:    input.Limit,
	}
	for _, r := range rows {
		out.EmployeeIDs = append(out.EmployeeIDs, &model.EmployeeShiftInfoBase{
			EmployeeId:          r.EmployeeID.Bytes,
			EmployeeName:        r.FullName,
			EmployeeCode:        r.EmployeeCode,
			EmployeeShiftName:   r.ShiftName,
			EmployeeShiftActive: r.CurrentShift,
		})
	}
	return out, nil
}

// GetListEmployeeInShift implements repository.IShiftUserRepository using sqlc.
func (s *ShiftUserRepository) GetListEmployeeInShift(ctx context.Context, input *model.GetListEmployyeShiftInput) (*model.GetListEmployyeShiftOutput, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}
	total, err := s.db.CountEmployeesInShiftCurrent(ctx, database.CountEmployeesInShiftCurrentParams{
		CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyID},
		ShiftID:   pgtype.UUID{Valid: true, Bytes: input.ShiftID},
	})
	if err != nil {
		return nil, err
	}
	rows, err := s.db.GetListEmployeeInShift(ctx, database.GetListEmployeeInShiftParams{
		CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyID},
		ShiftID:   pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		Limit:     input.Limit,
		Offset:    input.Offset,
	})
	if err != nil {
		return nil, err
	}
	out := &model.GetListEmployyeShiftOutput{
		EmployeeIDs: make([]*model.EmployeeShiftInfoBase, 0, len(rows)),
		Total:       int32(total),
		PageSize:    input.Limit,
	}
	for _, r := range rows {
		out.EmployeeIDs = append(out.EmployeeIDs, &model.EmployeeShiftInfoBase{
			EmployeeId:          r.EmployeeID.Bytes,
			EmployeeName:        r.FullName,
			EmployeeCode:        r.EmployeeCode,
			EmployeeShiftName:   r.ShiftName,
			EmployeeShiftActive: r.CurrentShift,
			ShiftEffectiveFrom:  r.ShiftEffectiveFrom.Time,
			ShiftEffectiveTo:    r.ShiftEffectiveTo.Time,
		})
	}
	return out, nil
}

// RemoveListShiftForEmployees implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) RemoveListShiftForEmployees(ctx context.Context, input *model.RemoveListShiftForEmployeesInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	listError := make([]error, 0)
	for _, employeeId := range input.EmployeeIDs {
		err := s.db.DeleteEmployeeShift(ctx, database.DeleteEmployeeShiftParams{
			EmployeeID: pgtype.UUID{Valid: true, Bytes: employeeId},
			ShiftID:    pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		})
		if err != nil {
			listError = append(listError, err)
		}
	}
	if len(listError) > 0 {
		return errors.New("one or more errors occurred while removing shifts for employees")
	}
	return nil
}

// IsUserManagetShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) IsUserManagetShift(ctx context.Context, input *model.IsUserManagetShiftInput) (bool, error) {
	_, err := s.db.IsUserManagetShift(
		ctx,
		database.IsUserManagetShiftParams{
			ShiftID:   pgtype.UUID{Valid: true, Bytes: input.ShiftID},
			CompanyID: pgtype.UUID{Valid: true, Bytes: input.CompanyUserID},
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

// DeleteEmployeeShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) DeleteEmployeeShift(ctx context.Context, input *model.DeleteEmployeeShiftInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	if err := s.db.DeleteEmployeeShift(
		ctx,
		database.DeleteEmployeeShiftParams{
			EmployeeID: pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
			ShiftID:    pgtype.UUID{Valid: true, Bytes: input.ShiftId},
		},
	); err != nil && err != pgx.ErrNoRows {
		return err
	}
	return nil
}

// DisableEmployeeShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) DisableEmployeeShift(ctx context.Context, input *model.DisableEmployeeShiftInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	if err := s.db.DisableEmployeeShift(
		ctx,
		database.DisableEmployeeShiftParams{
			EmployeeID: pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
			ShiftID:    pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		},
	); err != nil && err != pgx.ErrNoRows {
		return err
	}
	return nil
}

// EnableEmployeeShift implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) EnableEmployeeShift(ctx context.Context, input *model.EnableEmployeeShiftIInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	if err := s.db.EnableEmployeeShift(
		ctx,
		database.EnableEmployeeShiftParams{
			EmployeeID: pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
			ShiftID:    pgtype.UUID{Valid: true, Bytes: input.ShiftID},
		},
	); err != nil && err != pgx.ErrNoRows {
		return err
	}
	return nil
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
			EmployeeID:    input.EmployeeID,
			ShiftID:       r.ShiftID.Bytes,
			EffectiveFrom: fromPgDate(r.EffectiveFrom),
			EffectiveTo:   fromPgDate(r.EffectiveTo),
			IsActive:      fromPgBoolValue(r.IsActive),
		})
	}
	return out, nil
}

// GetShiftEmployeeAll returns all shift assignments for an employee with pagination (no date filter).
func (s *ShiftUserRepository) GetShiftEmployeeAll(ctx context.Context, input *model.GetShiftEmployeeAllInput) ([]*model.EmployeeShiftRow, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}
	rows, err := s.pool.Query(ctx, queryGetShiftEmployeeAll, input.EmployeeID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]*model.EmployeeShiftRow, 0)
	for rows.Next() {
		var (
			shiftID pgtype.UUID
			effFrom pgtype.Date
			effTo   pgtype.Date
			isAct   pgtype.Bool
		)
		if scanErr := rows.Scan(&shiftID, &effFrom, &effTo, &isAct); scanErr != nil {
			return nil, scanErr
		}
		out = append(out, &model.EmployeeShiftRow{
			EmployeeID:    input.EmployeeID,
			ShiftID:       shiftID.Bytes,
			EffectiveFrom: fromPgDate(effFrom),
			EffectiveTo:   fromPgDate(effTo),
			IsActive:      fromPgBoolValue(isAct),
		})
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	return out, nil
}

// EditEffectiveShiftForEmployee implements repository.IShiftUserRepository.
func (s *ShiftUserRepository) EditEffectiveShiftForEmployee(ctx context.Context, input *model.EditEffectiveShiftForEmployeeInput) error {
	if input == nil {
		return errors.New("input cannot be nil")
	}

	return s.db.EditEffectiveShiftForEmployee(ctx, database.EditEffectiveShiftForEmployeeParams{
		EmployeeID:    pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
		EffectiveFrom: toPgDate(input.EffectiveFrom),
		EffectiveTo:   toPgDate(input.EffectiveTo),
	})
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
func toPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func fromPgDate(d pgtype.Date) time.Time {
	if d.Valid {
		return d.Time
	}
	return time.Time{}
}

func fromPgBoolValue(b pgtype.Bool) bool {
	if b.Valid {
		return b.Bool
	}
	return false
}
