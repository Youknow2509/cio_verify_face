package repository

import (
	"context"
	"encoding/json"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/repository"
	db "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/gen"
)

/**
 * Struct impl IAuditRepository
 */
type AuditRepository struct {
	q db.Queries
}

// AddAuditLog implements repository.IAuditRepository.
func (a *AuditRepository) AddAuditLog(ctx context.Context, log *model.AuditLog) error {
	ipAddress, err := netip.ParseAddr(log.IpAddress)
	if err != nil {
		return err
	}
	oldValuesBytes, err := json.Marshal(log.OldValues)
	if err != nil {
		return err
	}
	newValuesBytes, err := json.Marshal(log.NewValues)
	if err != nil {
		return err
	}
	// Convert int64 timestamp to time.Time
	timestamp := time.Unix(log.Timestamp, 0)
	return a.q.AddAudit(ctx, db.AddAuditParams{
		UserID:       pgtype.UUID{Bytes: log.UserId, Valid: true},
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   pgtype.UUID{Bytes: log.ResourceId, Valid: true},
		OldValues:    oldValuesBytes,
		NewValues:    newValuesBytes,
		IpAddress:    &ipAddress,
		UserAgent:    pgtype.Text{String: log.UserAgent, Valid: true},
		Timestamp:    pgtype.Timestamptz{Time: timestamp, Valid: true},
	})
}

//

/**
 * New AuditRepository
 */
func NewAuditRepository(client *pgxpool.Pool) domainRepository.IAuditRepository {
	return &AuditRepository{
		q: *db.New(client),
	}
}
