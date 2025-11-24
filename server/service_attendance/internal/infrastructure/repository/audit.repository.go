package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
)

/**
 * Struct impl IAuditRepository
 */
type AuditRepository struct {
	dbSession *gocql.Session
}

// AddAuditLog implements repository.IAuditRepository.
func (a *AuditRepository) AddAuditLog(ctx context.Context, log *model.AuditLog) error {
	// 	INSERT INTO audit_logs (
	//     company_id, year_month, created_at, actor_id,
	//     action_category, action_name, resource_type, resource_id,
	//     details, ip_address, user_agent, status
	// ) VALUES (
	//     uuid_company, '2023-10', toTimestamp(now()), uuid_admin,
	//     'HR_MANAGEMENT', 'UPDATE_ATTENDANCE', 'ATTENDANCE_RECORD', 'rec_001',
	//     {'reason': 'Fixed forgot checkout', 'old_val': 'missing', 'new_val': '17:30'},
	//     '192.168.1.1', 'Mozilla/5.0...', 'SUCCESS'
	// );
	sql_raw := `INSERT INTO audit_logs (
		company_id, year_month, created_at, actor_id, 
		action_category, action_name, resource_type, resource_id,
		details, ip_address, user_agent, status
		) VALUES (
		?, ?, toTimestamp(now()), ?, 
		?, ?, ?, ?,
		?, ?, ?, ?
	);`
	return a.dbSession.Query(sql_raw,
		log.CompanyID,
		log.YearMonth,
		log.ActorID,
		log.ActionCategory,
		log.ActionName,
		log.ResourceType,
		log.ResourceID,
		log.Details,
		log.IP_Address,
		log.UserAgent,
		log.Status,
	).WithContext(ctx).Exec()
}

//

/**
 * New AuditRepository
 */
func NewAuditRepository(session *gocql.Session) domainRepository.IAuditRepository {
	return &AuditRepository{
		dbSession: session,
	}
}
