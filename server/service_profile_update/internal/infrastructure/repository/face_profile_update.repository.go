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
 * Struct impl IFaceProfileUpdateRequestRepository
 */
type FaceProfileUpdateRequestRepository struct {
	q db.Queries
}

// CreateRequest implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) CreateRequest(ctx context.Context, req *model.FaceProfileUpdateRequest) error {
	// Marshal metadata to JSON bytes
	metaDataBytes, err := json.Marshal(req.MetaData)
	if err != nil {
		return err
	}

	reasonText := pgtype.Text{}
	if req.Reason != nil {
		reasonText = pgtype.Text{String: *req.Reason, Valid: true}
	}

	return r.q.CreateFaceProfileUpdateRequest(ctx, db.CreateFaceProfileUpdateRequestParams{
		RequestID:           pgtype.UUID{Bytes: req.RequestID, Valid: true},
		UserID:              pgtype.UUID{Bytes: req.UserID, Valid: true},
		CompanyID:           pgtype.UUID{Bytes: req.CompanyID, Valid: true},
		Status:              pgtype.Int2{Int16: int16(req.Status), Valid: true},
		RequestMonth:        req.RequestMonth,
		RequestCountInMonth: pgtype.Int4{Int32: int32(req.RequestCountInMonth), Valid: true},
		Reason:              reasonText,
		MetaData:            metaDataBytes,
		CreatedAt:           pgtype.Timestamptz{Time: req.CreatedAt, Valid: true},
		UpdatedAt:           pgtype.Timestamptz{Time: req.UpdatedAt, Valid: true},
	})
}

// GetRequestByID implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) GetRequestByID(ctx context.Context, requestID, companyID uuid.UUID) (*model.FaceProfileUpdateRequest, error) {
	result, err := r.q.GetFaceProfileUpdateRequestByID(ctx, db.GetFaceProfileUpdateRequestByIDParams{
		RequestID: pgtype.UUID{Bytes: requestID, Valid: true},
		CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapToFaceProfileUpdateRequest(&result), nil
}

// GetRequestByToken implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) GetRequestByToken(ctx context.Context, token string) (*model.FaceProfileUpdateRequest, error) {
	result, err := r.q.GetFaceProfileUpdateRequestByToken(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapToFaceProfileUpdateRequest(&result), nil
}

// GetPendingRequestsByCompany implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) GetPendingRequestsByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*model.FaceProfileUpdateRequest, error) {
	results, err := r.q.GetPendingRequestsByCompany(ctx, db.GetPendingRequestsByCompanyParams{
		CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, err
	}

	requests := make([]*model.FaceProfileUpdateRequest, 0, len(results))
	for i := range results {
		requests = append(requests, mapToFaceProfileUpdateRequest(&results[i]))
	}

	return requests, nil
}

// GetRequestsByUserAndMonth implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) GetRequestsByUserAndMonth(ctx context.Context, userID uuid.UUID, month string) ([]*model.FaceProfileUpdateRequest, error) {
	results, err := r.q.GetRequestsByUserAndMonth(ctx, db.GetRequestsByUserAndMonthParams{
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		RequestMonth: month,
	})
	if err != nil {
		return nil, err
	}

	requests := make([]*model.FaceProfileUpdateRequest, 0, len(results))
	for i := range results {
		requests = append(requests, mapToFaceProfileUpdateRequest(&results[i]))
	}

	return requests, nil
}

// CountRequestsByUserInMonth implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) CountRequestsByUserInMonth(ctx context.Context, userID uuid.UUID, month string) (int, error) {
	count, err := r.q.CountRequestsByUserInMonth(ctx, db.CountRequestsByUserInMonthParams{
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		RequestMonth: month,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// UpdateRequestStatus implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) UpdateRequestStatus(ctx context.Context, requestID, companyID uuid.UUID, status model.RequestStatus) error {
	return r.q.UpdateRequestStatus(ctx, db.UpdateRequestStatusParams{
		RequestID: pgtype.UUID{Bytes: requestID, Valid: true},
		CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
		Status:    pgtype.Int2{Int16: int16(status), Valid: true},
	})
}

// ApproveRequest implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) ApproveRequest(ctx context.Context, requestID, companyID, approvedBy uuid.UUID, updateToken string, expiresAt interface{}) error {
	expiresTime, ok := expiresAt.(time.Time)
	if !ok {
		return errors.New("expiresAt must be time.Time")
	}

	return r.q.ApproveRequest(ctx, db.ApproveRequestParams{
		RequestID:           pgtype.UUID{Bytes: requestID, Valid: true},
		CompanyID:           pgtype.UUID{Bytes: companyID, Valid: true},
		ApprovedBy:          pgtype.UUID{Bytes: approvedBy, Valid: true},
		UpdateToken:         pgtype.Text{String: updateToken, Valid: true},
		UpdateLinkExpiresAt: pgtype.Timestamptz{Time: expiresTime, Valid: true},
	})
}

// RejectRequest implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) RejectRequest(ctx context.Context, requestID, companyID, rejectedBy uuid.UUID, reason string) error {
	return r.q.RejectRequest(ctx, db.RejectRequestParams{
		RequestID:       pgtype.UUID{Bytes: requestID, Valid: true},
		CompanyID:       pgtype.UUID{Bytes: companyID, Valid: true},
		ApprovedBy:      pgtype.UUID{Bytes: rejectedBy, Valid: true},
		RejectionReason: pgtype.Text{String: reason, Valid: true},
	})
}

// CompleteRequest implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) CompleteRequest(ctx context.Context, requestID, companyID uuid.UUID) error {
	return r.q.CompleteRequest(ctx, db.CompleteRequestParams{
		RequestID: pgtype.UUID{Bytes: requestID, Valid: true},
		CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
	})
}

// MarkExpiredRequests implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) MarkExpiredRequests(ctx context.Context) (int64, error) {
	return r.q.MarkExpiredRequests(ctx)
}

// HasPendingRequest implements repository.IFaceProfileUpdateRequestRepository
func (r *FaceProfileUpdateRequestRepository) HasPendingRequest(ctx context.Context, userID, companyID uuid.UUID) (bool, error) {
	result, err := r.q.HasPendingRequest(ctx, db.HasPendingRequestParams{
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

// Helper function to map DB result to domain model (works with FaceProfileUpdateRequest directly)
func mapToFaceProfileUpdateRequest(r *db.FaceProfileUpdateRequest) *model.FaceProfileUpdateRequest {
	// Unmarshal metadata
	var metaData map[string]interface{}
	if len(r.MetaData) > 0 {
		_ = json.Unmarshal(r.MetaData, &metaData)
	}
	if metaData == nil {
		metaData = make(map[string]interface{})
	}

	req := &model.FaceProfileUpdateRequest{
		RequestID:           r.RequestID.Bytes,
		UserID:              r.UserID.Bytes,
		CompanyID:           r.CompanyID.Bytes,
		Status:              model.RequestStatus(r.Status.Int16),
		RequestMonth:        r.RequestMonth,
		RequestCountInMonth: int(r.RequestCountInMonth.Int32),
		MetaData:            metaData,
		CreatedAt:           r.CreatedAt.Time,
		UpdatedAt:           r.UpdatedAt.Time,
	}

	if r.Reason.Valid {
		reason := r.Reason.String
		req.Reason = &reason
	}
	if r.UpdateToken.Valid {
		updateToken := r.UpdateToken.String
		req.UpdateToken = &updateToken
	}
	if r.UpdateLinkExpiresAt.Valid {
		req.UpdateLinkExpiresAt = &r.UpdateLinkExpiresAt.Time
	}
	if r.ApprovedBy.Valid {
		approvedBy := uuid.UUID(r.ApprovedBy.Bytes)
		req.ApprovedBy = &approvedBy
	}
	if r.ApprovedAt.Valid {
		req.ApprovedAt = &r.ApprovedAt.Time
	}
	if r.RejectionReason.Valid {
		rejectionReason := r.RejectionReason.String
		req.RejectionReason = &rejectionReason
	}

	return req
}

/**
 * NewFaceProfileUpdateRequestRepository creates a new repository
 */
func NewFaceProfileUpdateRequestRepository(client *pgxpool.Pool) domainRepository.IFaceProfileUpdateRequestRepository {
	return &FaceProfileUpdateRequestRepository{
		q: *db.New(client),
	}
}
