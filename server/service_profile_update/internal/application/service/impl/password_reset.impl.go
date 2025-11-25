package impl

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/errors"
	appModel "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/cache"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/model"
	domainMQ "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/mq"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/repository"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
)

// PasswordResetServiceImpl implements IPasswordResetService
type PasswordResetServiceImpl struct{}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService() *PasswordResetServiceImpl {
	return &PasswordResetServiceImpl{}
}

// ResetEmployeePassword - Manager resets an employee's password
func (s *PasswordResetServiceImpl) ResetEmployeePassword(ctx context.Context, input *appModel.ResetEmployeePasswordInput) (*appModel.ResetEmployeePasswordOutput, *appErrors.Error) {
	if input.Session == nil {
		return nil, appErrors.ErrUnauthorized.WithDetails("session required")
	}

	// Authorization check - must be company admin or system admin
	if err := s.checkAuthorization(input.Session, nil, domainModel.RoleCompanyAdmin, false); err != nil {
		return nil, err
	}

	managerID, err := uuid.Parse(input.Session.UserID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid manager ID")
	}

	employeeID, err := uuid.Parse(input.EmployeeID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid employee ID")
	}

	companyID, err := uuid.Parse(input.Session.CompanyID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid company ID")
	}

	// Get repositories
	userRepo, err := domainRepo.GetUserRepository()
	if err != nil {
		global.Logger.Error("Failed to get user repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	prrRepo, err := domainRepo.GetPasswordResetRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get password reset request repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	// Check spam - manager cannot spam reset for same employee
	if blocked := s.checkPasswordResetSpam(ctx, managerID, employeeID); blocked {
		return nil, appErrors.ErrPasswordResetSpam.WithDetails("please wait before resetting this employee's password again")
	}

	// Check if manager is allowed to reset this employee's password
	if domainModel.Role(input.Session.Role) != domainModel.RoleSystemAdmin {
		// Check if employee belongs to manager's company
		belongs, err := userRepo.UserBelongsToCompany(ctx, employeeID, companyID)
		if err != nil {
			global.Logger.Error("Failed to check employee company", err)
			return nil, appErrors.ErrServiceUnavailable
		}
		if !belongs {
			return nil, appErrors.ErrForbidden.WithDetails("employee does not belong to your company")
		}
	}

	// Get employee info
	employeeInfo, err := userRepo.GetUserByID(ctx, employeeID)
	if err != nil {
		global.Logger.Error("Failed to get employee info", err)
		return nil, appErrors.ErrServiceUnavailable
	}
	if employeeInfo == nil {
		return nil, appErrors.ErrEmployeeNotFound
	}

	// Check hourly limit for manager
	if exceeded := s.checkManagerHourlyLimit(ctx, prrRepo, managerID); exceeded {
		return nil, appErrors.ErrPasswordResetSpam.WithDetails("hourly password reset limit exceeded")
	}

	// Generate new random password
	newPassword := s.generateRandomPassword(12)
	salt := s.generateSalt()
	passwordHash := s.hashPassword(newPassword, salt)

	// Update password in database
	if err := userRepo.UpdateUserPassword(ctx, employeeID, salt, passwordHash); err != nil {
		global.Logger.Error("Failed to update user password", err)
		return nil, appErrors.ErrPasswordResetFailed.WithDetails("failed to update password")
	}

	// Create password reset request record
	requestID := uuid.New()
	now := time.Now()
	resetRequest := &domainModel.PasswordResetRequest{
		RequestID:   requestID,
		UserID:      employeeID,
		CompanyID:   &companyID,
		RequestedBy: managerID,
		Status:      domainModel.PasswordResetStatusPending,
		MetaData: map[string]interface{}{
			"client_ip":      input.Session.ClientIP,
			"user_agent":     input.Session.ClientAgent,
			"employee_email": employeeInfo.Email,
			"employee_name":  employeeInfo.FullName,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := prrRepo.CreateRequest(ctx, resetRequest); err != nil {
		global.Logger.Error("Failed to create password reset request record", err)
		// Don't fail - password is already reset
	}

	// Send notification via Kafka
	kafkaMessageID, kafkaErr := s.sendPasswordResetNotification(ctx, employeeInfo, newPassword, requestID)
	if kafkaErr != nil {
		global.Logger.Error("Failed to send password reset notification", kafkaErr)
		// Update request status to failed
		_ = prrRepo.UpdateRequestStatus(ctx, requestID, domainModel.PasswordResetStatusFailed, "")
	} else {
		// Update request status to sent
		_ = prrRepo.UpdateRequestStatus(ctx, requestID, domainModel.PasswordResetStatusSent, kafkaMessageID)
	}

	// Set spam prevention marker
	s.setPasswordResetSpamMarker(ctx, managerID, employeeID)

	return &appModel.ResetEmployeePasswordOutput{
		Success:     true,
		Message:     "Password reset successfully. New password sent to employee's email.",
		NewPassword: newPassword, // Only returned to manager for reference
	}, nil
}

// =================================
// Helper Methods:
// =================================

func (s *PasswordResetServiceImpl) checkAuthorization(session *appModel.SessionInfo, requestedCompanyID *string, minRole domainModel.Role, allowEmployeeSelfAccess bool) *appErrors.Error {
	if session == nil {
		return appErrors.ErrUnauthorized.WithDetails("session info required")
	}

	userRole := domainModel.Role(session.Role)

	// System admin has full access
	if userRole == domainModel.RoleSystemAdmin {
		return nil
	}

	// Check company match for non-system-admin users
	if requestedCompanyID != nil && session.CompanyID != *requestedCompanyID {
		return appErrors.ErrForbidden.WithDetails("access denied: you can only access your own company data")
	}

	// Check role
	if userRole > minRole {
		if !allowEmployeeSelfAccess || userRole != domainModel.RoleEmployee {
			return appErrors.ErrForbidden.WithDetails("access denied: insufficient permissions")
		}
	}

	return nil
}

func (s *PasswordResetServiceImpl) checkPasswordResetSpam(ctx context.Context, managerID, employeeID uuid.UUID) bool {
	key := constants.CacheKeyPrefixPasswordResetSpam + managerID.String() + ":" + employeeID.String()

	// Check local cache first
	if localCache, err := domainCache.GetLocalCache(); err == nil {
		if exists, _ := localCache.Exists(ctx, key); exists {
			return true
		}
	}

	// Check distributed cache
	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		if exists, _ := distCache.Exists(ctx, key); exists {
			return true
		}
	}

	return false
}

func (s *PasswordResetServiceImpl) setPasswordResetSpamMarker(ctx context.Context, managerID, employeeID uuid.UUID) {
	key := constants.CacheKeyPrefixPasswordResetSpam + managerID.String() + ":" + employeeID.String()
	ttl := int64(constants.PasswordResetCooldownSeconds)

	if localCache, err := domainCache.GetLocalCache(); err == nil {
		_ = localCache.SetTTL(ctx, key, "1", ttl)
	}
	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.SetTTL(ctx, key, "1", ttl)
	}
}

func (s *PasswordResetServiceImpl) checkManagerHourlyLimit(ctx context.Context, repo domainRepo.IPasswordResetRequestRepository, managerID uuid.UUID) bool {
	since := time.Now().Add(-1 * time.Hour)
	count, err := repo.CountRequestsByManagerInWindow(ctx, managerID, since)
	if err != nil {
		global.Logger.Error("Failed to count manager password reset requests", err)
		return false // Allow if count fails
	}
	return count >= constants.MaxPasswordResetsPerManagerPerHour
}

func (s *PasswordResetServiceImpl) generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to simple random if crypto/rand fails
			password[i] = charset[i%len(charset)]
		} else {
			password[i] = charset[n.Int64()]
		}
	}
	return string(password)
}

func (s *PasswordResetServiceImpl) generateSalt() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return uuid.New().String()[:16]
	}
	return fmt.Sprintf("%x", bytes)
}

func (s *PasswordResetServiceImpl) hashPassword(password, salt string) string {
	// Use the same SHA-256 hashing as the existing auth service
	// Hash: SHA256(password + salt)
	saltedPassword := password + salt
	hash := sha256.Sum256([]byte(saltedPassword))
	return hex.EncodeToString(hash[:])
}

func (s *PasswordResetServiceImpl) sendPasswordResetNotification(ctx context.Context, employee *domainModel.UserInfo, newPassword string, requestID uuid.UUID) (string, error) {
	kafkaWriter, err := domainMQ.GetKafkaWriter()
	if err != nil {
		return "", err
	}

	messageID := uuid.New().String()

	payload := map[string]interface{}{
		"message_id":   messageID,
		"request_id":   requestID.String(),
		"event_type":   constants.KafkaEventTypePasswordReset,
		"user_id":      employee.UserID.String(),
		"email":        employee.Email,
		"full_name":    employee.FullName,
		"new_password": newPassword,
		"created_at":   time.Now().UTC().Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal kafka payload: %w", err)
	}

	// Use RequireAck to ensure message is received
	if err := kafkaWriter.WriteMessageRequireAck(ctx, constants.KafkaTopicPasswordResetNotifications, employee.UserID.String(), payloadBytes); err != nil {
		return "", fmt.Errorf("failed to send kafka message: %w", err)
	}

	return messageID, nil
}
