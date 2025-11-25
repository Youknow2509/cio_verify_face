package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/repository"
	db "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/gen"
)

/**
 * Struct impl IPasswordResetRequestRepository
 */
type PasswordResetRequestRepository struct {
	q db.Queries
}

// CreateRequest implements repository.IPasswordResetRequestRepository
func (r *PasswordResetRequestRepository) CreateRequest(ctx context.Context, req *model.PasswordResetRequest) error {
	companyID := pgtype.UUID{}
	if req.CompanyID != nil {
		companyID = pgtype.UUID{Bytes: *req.CompanyID, Valid: true}
	}

	// Marshal metadata to JSON bytes
	metaDataBytes, err := json.Marshal(req.MetaData)
	if err != nil {
		return err
	}

	return r.q.CreatePasswordResetRequest(ctx, db.CreatePasswordResetRequestParams{
		RequestID:   pgtype.UUID{Bytes: req.RequestID, Valid: true},
		UserID:      pgtype.UUID{Bytes: req.UserID, Valid: true},
		CompanyID:   companyID,
		RequestedBy: pgtype.UUID{Bytes: req.RequestedBy, Valid: true},
		Status:      pgtype.Int2{Int16: int16(req.Status), Valid: true},
		MetaData:    metaDataBytes,
		CreatedAt:   pgtype.Timestamptz{Time: req.CreatedAt, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: req.UpdatedAt, Valid: true},
	})
}

// GetRequestByID implements repository.IPasswordResetRequestRepository
func (r *PasswordResetRequestRepository) GetRequestByID(ctx context.Context, requestID uuid.UUID) (*model.PasswordResetRequest, error) {
	result, err := r.q.GetPasswordResetRequestByID(ctx, pgtype.UUID{Bytes: requestID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapToPasswordResetRequest(&result), nil
}

// GetRecentRequestsByManagerForUser implements repository.IPasswordResetRequestRepository
func (r *PasswordResetRequestRepository) GetRecentRequestsByManagerForUser(ctx context.Context, managerID, userID uuid.UUID, since interface{}) ([]*model.PasswordResetRequest, error) {
	sinceTime, ok := since.(time.Time)
	if !ok {
		return nil, errors.New("since must be time.Time")
	}

	results, err := r.q.GetRecentRequestsByManagerForUser(ctx, db.GetRecentRequestsByManagerForUserParams{
		RequestedBy: pgtype.UUID{Bytes: managerID, Valid: true},
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: sinceTime, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	requests := make([]*model.PasswordResetRequest, 0, len(results))
	for i := range results {
		requests = append(requests, mapToPasswordResetRequest(&results[i]))
	}

	return requests, nil
}

// UpdateRequestStatus implements repository.IPasswordResetRequestRepository
func (r *PasswordResetRequestRepository) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status model.PasswordResetStatus, kafkaMessageID string) error {
	return r.q.UpdatePasswordResetStatus(ctx, db.UpdatePasswordResetStatusParams{
		RequestID:      pgtype.UUID{Bytes: requestID, Valid: true},
		Status:         pgtype.Int2{Int16: int16(status), Valid: true},
		KafkaMessageID: pgtype.Text{String: kafkaMessageID, Valid: kafkaMessageID != ""},
	})
}

// CountRequestsByManagerInWindow implements repository.IPasswordResetRequestRepository
func (r *PasswordResetRequestRepository) CountRequestsByManagerInWindow(ctx context.Context, managerID uuid.UUID, since interface{}) (int, error) {
	sinceTime, ok := since.(time.Time)
	if !ok {
		return 0, errors.New("since must be time.Time")
	}

	count, err := r.q.CountRequestsByManagerInWindow(ctx, db.CountRequestsByManagerInWindowParams{
		RequestedBy: pgtype.UUID{Bytes: managerID, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: sinceTime, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Helper function to map DB result to domain model
func mapToPasswordResetRequest(r *db.PasswordResetRequest) *model.PasswordResetRequest {
	// Unmarshal metadata
	var metaData map[string]interface{}
	if len(r.MetaData) > 0 {
		_ = json.Unmarshal(r.MetaData, &metaData)
	}
	if metaData == nil {
		metaData = make(map[string]interface{})
	}

	req := &model.PasswordResetRequest{
		RequestID:   r.RequestID.Bytes,
		UserID:      r.UserID.Bytes,
		RequestedBy: r.RequestedBy.Bytes,
		Status:      model.PasswordResetStatus(r.Status.Int16),
		MetaData:    metaData,
		CreatedAt:   r.CreatedAt.Time,
		UpdatedAt:   r.UpdatedAt.Time,
	}

	if r.CompanyID.Valid {
		companyID := uuid.UUID(r.CompanyID.Bytes)
		req.CompanyID = &companyID
	}
	if r.KafkaMessageID.Valid {
		kafkaID := r.KafkaMessageID.String
		req.KafkaMessageID = &kafkaID
	}
	if r.KafkaSentAt.Valid {
		req.KafkaSentAt = &r.KafkaSentAt.Time
	}

	return req
}

/**
 * NewPasswordResetRequestRepository creates a new repository
 */
func NewPasswordResetRequestRepository(client *pgxpool.Pool) domainRepository.IPasswordResetRequestRepository {
	return &PasswordResetRequestRepository{
		q: *db.New(client),
	}
}
