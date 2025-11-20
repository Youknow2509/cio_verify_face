package repository

import (
	"context"
	"errors"
	"time"

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

// GetListTimeShiftEmployee implements repository.IUserRepository.
func (u *UserRepository) GetListTimeShiftEmployee(ctx context.Context, input *domainModel.GetListTimeShiftEmployeeInput) ([]domainModel.ShiftTimeEmployee, error) {
	reps, err := u.q.GetListTimeShiftEmployee(
		ctx,
		db.GetListTimeShiftEmployeeParams{
			CompanyID:  pgtype.UUID{Valid: true, Bytes: input.CompanyID},
			EmployeeID: pgtype.UUID{Valid: true, Bytes: input.EmployeeID},
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domainModel.ShiftTimeEmployee{}, nil
		}
		return nil, err
	}

	result := make([]domainModel.ShiftTimeEmployee, len(reps))
	for i, r := range reps {
		var effectiveTo *time.Time
		if r.EffectiveTo.Valid {
			t := r.EffectiveTo.Time
			effectiveTo = &t
		}

		// Convert pgtype.Time to time.Time (using date 0000-01-01 as base)
		startTime := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(r.StartTime.Microseconds) * time.Microsecond)
		endTime := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(r.EndTime.Microseconds) * time.Microsecond)

		result[i] = domainModel.ShiftTimeEmployee{
			StartTime:             startTime,
			EndTime:               endTime,
			GracePeriodMinutes:    int(r.GracePeriodMinutes.Int32),
			EarlyDepartureMinutes: int(r.EarlyDepartureMinutes.Int32),
			WorkDays:              r.WorkDays,
			EffectiveFrom:         r.EffectiveFrom.Time,
			EffectiveTo:           effectiveTo,
		}
	}
	return result, nil
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
