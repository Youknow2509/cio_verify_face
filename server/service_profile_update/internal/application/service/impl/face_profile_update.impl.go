package impl

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/errors"
	appModel "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/cache"
	domainGrpc "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/grpc"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/repository"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
)

// FaceProfileUpdateServiceImpl implements IFaceProfileUpdateService
type FaceProfileUpdateServiceImpl struct {
	baseURL string
}

// NewFaceProfileUpdateService creates a new face profile update service
func NewFaceProfileUpdateService(baseURL string) *FaceProfileUpdateServiceImpl {
	return &FaceProfileUpdateServiceImpl{
		baseURL: baseURL,
	}
}

// CreateUpdateRequest - Employee creates a new face profile update request
func (s *FaceProfileUpdateServiceImpl) CreateUpdateRequest(ctx context.Context, input *appModel.CreateFaceProfileUpdateRequestInput) (*appModel.CreateFaceProfileUpdateRequestOutput, *appErrors.Error) {
	if input.Session == nil {
		return nil, appErrors.ErrUnauthorized.WithDetails("session required")
	}

	userID, err := uuid.Parse(input.Session.UserID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid user ID")
	}

	companyID, err := uuid.Parse(input.Session.CompanyID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid company ID")
	}

	// Get repositories
	fprRepo, err := domainRepo.GetFaceProfileUpdateRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get face profile update request repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	// Check for spam using local cache first
	if blocked := s.checkSpamLocal(ctx, userID); blocked {
		return nil, appErrors.ErrSpamDetected.WithDetails("please wait before creating another request")
	}

	// Check for spam using distributed cache
	if blocked := s.checkSpamDistributed(ctx, userID); blocked {
		return nil, appErrors.ErrSpamDetected.WithDetails("please wait before creating another request")
	}

	// Check if user already has a pending request
	hasPending, err := s.checkPendingRequest(ctx, fprRepo, userID, companyID)
	if err != nil {
		global.Logger.Error("Failed to check pending request", err)
		return nil, appErrors.ErrServiceUnavailable
	}
	if hasPending {
		return nil, appErrors.ErrRequestAlreadyPending
	}

	// Check monthly limit
	currentMonth := time.Now().Format("2006-01")
	count, err := s.getMonthlyRequestCount(ctx, fprRepo, userID, currentMonth)
	if err != nil {
		global.Logger.Error("Failed to get monthly request count", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	if count >= constants.MaxFaceProfileRequestsPerMonth {
		return &appModel.CreateFaceProfileUpdateRequestOutput{
			Status:            "rejected",
			Message:           fmt.Sprintf("Monthly limit of %d requests reached. Please contact your manager.", constants.MaxFaceProfileRequestsPerMonth),
			CanRetry:          false,
			RequestsRemaining: 0,
		}, nil
	}

	// Create the request
	requestID := uuid.New()
	now := time.Now()
	reason := input.Reason
	request := &domainModel.FaceProfileUpdateRequest{
		RequestID:           requestID,
		UserID:              userID,
		CompanyID:           companyID,
		Status:              domainModel.RequestStatusPending,
		RequestMonth:        currentMonth,
		RequestCountInMonth: count + 1,
		Reason:              &reason,
		MetaData:            make(map[string]interface{}),
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Add metadata
	request.MetaData["client_ip"] = input.Session.ClientIP
	request.MetaData["user_agent"] = input.Session.ClientAgent

	if err := fprRepo.CreateRequest(ctx, request); err != nil {
		global.Logger.Error("Failed to create face profile update request", err)
		return nil, appErrors.ErrServiceUnavailable.WithDetails("failed to create request")
	}

	// Set spam prevention marker
	s.setSpamMarker(ctx, userID)

	// Invalidate caches
	s.invalidatePendingCache(ctx, userID, companyID)
	s.invalidateMonthlyCountCache(ctx, userID, currentMonth)

	return &appModel.CreateFaceProfileUpdateRequestOutput{
		RequestID:         requestID.String(),
		Status:            "pending",
		Message:           "Request submitted successfully. Please wait for manager approval.",
		CanRetry:          true,
		RequestsRemaining: constants.MaxFaceProfileRequestsPerMonth - count - 1,
	}, nil
}

// GetMyUpdateRequests - Employee gets their own update requests
func (s *FaceProfileUpdateServiceImpl) GetMyUpdateRequests(ctx context.Context, input *appModel.GetMyUpdateRequestsInput) (*appModel.GetMyUpdateRequestsOutput, *appErrors.Error) {
	if input.Session == nil {
		return nil, appErrors.ErrUnauthorized.WithDetails("session required")
	}

	userID, err := uuid.Parse(input.Session.UserID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid user ID")
	}

	fprRepo, err := domainRepo.GetFaceProfileUpdateRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	month := input.Month
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	requests, err := fprRepo.GetRequestsByUserAndMonth(ctx, userID, month)
	if err != nil {
		global.Logger.Error("Failed to get requests", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	dtos := make([]*appModel.FaceProfileUpdateRequestDTO, 0, len(requests))
	for _, req := range requests {
		dtos = append(dtos, appModel.ToFaceProfileUpdateRequestDTO(req, s.baseURL))
	}

	return &appModel.GetMyUpdateRequestsOutput{
		Requests: dtos,
		Total:    len(dtos),
	}, nil
}

// GetPendingRequests - Manager gets pending requests for their company
func (s *FaceProfileUpdateServiceImpl) GetPendingRequests(ctx context.Context, input *appModel.GetPendingRequestsInput) (*appModel.GetPendingRequestsOutput, *appErrors.Error) {
	if input.Session == nil {
		return nil, appErrors.ErrUnauthorized.WithDetails("session required")
	}

	// Authorization check - must be company admin or system admin
	if err := s.checkAuthorization(input.Session, nil, domainModel.RoleCompanyAdmin, false); err != nil {
		return nil, err
	}

	companyID, err := uuid.Parse(input.Session.CompanyID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid company ID")
	}

	fprRepo, err := domainRepo.GetFaceProfileUpdateRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	limit := input.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	requests, err := fprRepo.GetPendingRequestsByCompany(ctx, companyID, limit, offset)
	if err != nil {
		global.Logger.Error("Failed to get pending requests", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	dtos := make([]*appModel.FaceProfileUpdateRequestDTO, 0, len(requests))
	for _, req := range requests {
		dtos = append(dtos, appModel.ToFaceProfileUpdateRequestDTO(req, s.baseURL))
	}

	return &appModel.GetPendingRequestsOutput{
		Requests: dtos,
		Total:    len(dtos),
	}, nil
}

// ApproveRequest - Manager approves a face profile update request
func (s *FaceProfileUpdateServiceImpl) ApproveRequest(ctx context.Context, input *appModel.ApproveRequestInput) (*appModel.ApproveRequestOutput, *appErrors.Error) {
	if input.Session == nil {
		return nil, appErrors.ErrUnauthorized.WithDetails("session required")
	}

	// Authorization check
	targetCompanyID := input.CompanyID
	if targetCompanyID == "" {
		targetCompanyID = input.Session.CompanyID
	}
	if err := s.checkAuthorization(input.Session, &targetCompanyID, domainModel.RoleCompanyAdmin, false); err != nil {
		return nil, err
	}

	requestID, err := uuid.Parse(input.RequestID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid request ID")
	}

	companyID, err := uuid.Parse(targetCompanyID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid company ID")
	}

	approverID, err := uuid.Parse(input.Session.UserID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid approver ID")
	}

	// Check for duplicate approval using distributed lock
	if !s.acquireApprovalLock(ctx, requestID) {
		return nil, appErrors.ErrConflict.WithDetails("request is being processed")
	}
	defer s.releaseApprovalLock(ctx, requestID)

	fprRepo, err := domainRepo.GetFaceProfileUpdateRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	// Get the request
	request, err := fprRepo.GetRequestByID(ctx, requestID, companyID)
	if err != nil {
		global.Logger.Error("Failed to get request", err)
		return nil, appErrors.ErrServiceUnavailable
	}
	if request == nil {
		return nil, appErrors.ErrRequestNotFound
	}

	// Check if already processed
	if request.Status != domainModel.RequestStatusPending {
		return nil, appErrors.ErrRequestAlreadyProcessed.WithDetails(fmt.Sprintf("request status: %d", request.Status))
	}

	// Generate update token
	updateToken := s.generateSecureToken()
	expiresAt := time.Now().Add(time.Duration(constants.UpdateLinkTTLSeconds) * time.Second)

	// Approve the request
	if err := fprRepo.ApproveRequest(ctx, requestID, companyID, approverID, updateToken, expiresAt); err != nil {
		global.Logger.Error("Failed to approve request", err)
		return nil, appErrors.ErrServiceUnavailable.WithDetails("failed to approve request")
	}

	// Cache the token for quick validation
	s.cacheUpdateToken(ctx, updateToken, request.UserID.String(), companyID.String(), expiresAt)

	updateLink := fmt.Sprintf("%s/api/v1/profile-update/face?token=%s", s.baseURL, updateToken)

	return &appModel.ApproveRequestOutput{
		RequestID:  requestID.String(),
		UpdateLink: updateLink,
		ExpiresAt:  expiresAt,
		Message:    "Request approved. Update link sent to employee.",
	}, nil
}

// RejectRequest - Manager rejects a face profile update request
func (s *FaceProfileUpdateServiceImpl) RejectRequest(ctx context.Context, input *appModel.RejectRequestInput) (*appModel.RejectRequestOutput, *appErrors.Error) {
	if input.Session == nil {
		return nil, appErrors.ErrUnauthorized.WithDetails("session required")
	}

	// Authorization check
	targetCompanyID := input.CompanyID
	if targetCompanyID == "" {
		targetCompanyID = input.Session.CompanyID
	}
	if err := s.checkAuthorization(input.Session, &targetCompanyID, domainModel.RoleCompanyAdmin, false); err != nil {
		return nil, err
	}

	requestID, err := uuid.Parse(input.RequestID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid request ID")
	}

	companyID, err := uuid.Parse(targetCompanyID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid company ID")
	}

	rejecterID, err := uuid.Parse(input.Session.UserID)
	if err != nil {
		return nil, appErrors.ErrInvalidInput.WithDetails("invalid rejecter ID")
	}

	fprRepo, err := domainRepo.GetFaceProfileUpdateRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	// Get the request
	request, err := fprRepo.GetRequestByID(ctx, requestID, companyID)
	if err != nil {
		global.Logger.Error("Failed to get request", err)
		return nil, appErrors.ErrServiceUnavailable
	}
	if request == nil {
		return nil, appErrors.ErrRequestNotFound
	}

	// Check if already processed
	if request.Status != domainModel.RequestStatusPending {
		return nil, appErrors.ErrRequestAlreadyProcessed.WithDetails(fmt.Sprintf("request status: %d", request.Status))
	}

	// Reject the request
	if err := fprRepo.RejectRequest(ctx, requestID, companyID, rejecterID, input.Reason); err != nil {
		global.Logger.Error("Failed to reject request", err)
		return nil, appErrors.ErrServiceUnavailable.WithDetails("failed to reject request")
	}

	return &appModel.RejectRequestOutput{
		RequestID: requestID.String(),
		Status:    "rejected",
		Message:   "Request rejected successfully.",
	}, nil
}

// ValidateUpdateToken - Validate update token before allowing face update
func (s *FaceProfileUpdateServiceImpl) ValidateUpdateToken(ctx context.Context, input *appModel.ValidateUpdateTokenInput) (*appModel.ValidateUpdateTokenOutput, *appErrors.Error) {
	if input.Token == "" {
		return nil, appErrors.ErrInvalidInput.WithDetails("token required")
	}

	// Try to get from cache first
	userID, companyID, expiresAt, found := s.getTokenFromCache(ctx, input.Token)
	if found {
		if time.Now().After(expiresAt) {
			return &appModel.ValidateUpdateTokenOutput{
				Valid:   false,
				Message: "Update token has expired",
			}, nil
		}
		return &appModel.ValidateUpdateTokenOutput{
			Valid:     true,
			UserID:    userID,
			CompanyID: companyID,
			ExpiresAt: expiresAt,
			Message:   "Token is valid",
		}, nil
	}

	// Fall back to database
	fprRepo, err := domainRepo.GetFaceProfileUpdateRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	request, err := fprRepo.GetRequestByToken(ctx, input.Token)
	if err != nil {
		global.Logger.Error("Failed to get request by token", err)
		return nil, appErrors.ErrServiceUnavailable
	}
	if request == nil {
		return &appModel.ValidateUpdateTokenOutput{
			Valid:   false,
			Message: "Invalid update token",
		}, nil
	}

	if request.Status != domainModel.RequestStatusApproved {
		return &appModel.ValidateUpdateTokenOutput{
			Valid:   false,
			Message: "Update token is no longer valid",
		}, nil
	}

	if request.UpdateLinkExpiresAt != nil && time.Now().After(*request.UpdateLinkExpiresAt) {
		return &appModel.ValidateUpdateTokenOutput{
			Valid:   false,
			Message: "Update token has expired",
		}, nil
	}

	// Cache the token for future lookups
	if request.UpdateLinkExpiresAt != nil {
		s.cacheUpdateToken(ctx, input.Token, request.UserID.String(), request.CompanyID.String(), *request.UpdateLinkExpiresAt)
	}

	return &appModel.ValidateUpdateTokenOutput{
		Valid:     true,
		UserID:    request.UserID.String(),
		CompanyID: request.CompanyID.String(),
		ExpiresAt: *request.UpdateLinkExpiresAt,
		Message:   "Token is valid",
	}, nil
}

// UpdateFaceProfile - Employee updates their face profile using valid token
func (s *FaceProfileUpdateServiceImpl) UpdateFaceProfile(ctx context.Context, input *appModel.UpdateFaceProfileInput) (*appModel.UpdateFaceProfileOutput, *appErrors.Error) {
	// First validate the token
	validateResult, appErr := s.ValidateUpdateToken(ctx, &appModel.ValidateUpdateTokenInput{Token: input.Token})
	if appErr != nil {
		return nil, appErr
	}
	if !validateResult.Valid {
		return nil, appErrors.ErrInvalidUpdateToken.WithDetails(validateResult.Message)
	}

	// Get the request to mark it as completed
	fprRepo, err := domainRepo.GetFaceProfileUpdateRequestRepository()
	if err != nil {
		global.Logger.Error("Failed to get repository", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	request, err := fprRepo.GetRequestByToken(ctx, input.Token)
	if err != nil || request == nil {
		global.Logger.Error("Failed to get request by token", err)
		return nil, appErrors.ErrServiceUnavailable
	}

	// Get face service client
	faceClient := domainGrpc.GetFaceServiceClient()
	if faceClient == nil {
		global.Logger.Error("Face service client not initialized", nil)
		return nil, appErrors.ErrServiceUnavailable.WithDetails("face service unavailable")
	}

	// Step 1: Get existing face profiles for the user
	getProfilesReq := &domainGrpc.GetUserProfilesRequest{
		UserID:     request.UserID.String(),
		CompanyID:  request.CompanyID.String(),
		PageNumber: 1,
		PageSize:   100,
	}

	existingProfiles, err := faceClient.GetUserProfiles(ctx, getProfilesReq)
	if err != nil {
		// Differentiate between 'not found' (acceptable) and actual errors
		errMsg := err.Error()
		if !strings.Contains(errMsg, "not found") && !strings.Contains(errMsg, "no profiles") {
			global.Logger.Error("Failed to get existing face profiles", err)
			return nil, appErrors.ErrServiceUnavailable.WithDetails("failed to communicate with face service")
		}
		// User has no existing profiles - this is acceptable for first-time enrollment
		existingProfiles = &domainGrpc.GetUserProfilesResponse{Profiles: []*domainGrpc.FaceProfile{}}
	}

	// Step 2: Delete all existing face profiles
	for _, profile := range existingProfiles.Profiles {
		if profile.DeletedAt == nil { // Only delete non-deleted profiles
			deleteReq := &domainGrpc.DeleteProfileRequest{
				ProfileID:  profile.ProfileID,
				CompanyID:  request.CompanyID.String(),
				HardDelete: false, // Soft delete
			}
			if _, delErr := faceClient.DeleteProfile(ctx, deleteReq); delErr != nil {
				global.Logger.Error("Failed to delete existing profile", delErr)
				// Continue to enrollment even if delete fails
			}
		}
	}

	// Step 3: Enroll the new face profile
	enrollReq := &domainGrpc.EnrollFaceRequest{
		ImageData:   input.ImageData,
		UserID:      request.UserID.String(),
		CompanyID:   request.CompanyID.String(),
		DeviceID:    "", // Not from device
		MakePrimary: true,
		Filename:    input.Filename,
	}

	enrollResp, err := faceClient.EnrollFace(ctx, enrollReq)
	if err != nil {
		global.Logger.Error("Failed to enroll new face profile", err)
		return nil, appErrors.ErrServiceUnavailable.WithDetails("failed to enroll face profile")
	}

	if !strings.EqualFold(enrollResp.Status, domainGrpc.FaceServiceStatusSuccess) {
		global.Logger.Error("Face enrollment failed", fmt.Errorf("%s", enrollResp.Message))
		return nil, appErrors.ErrServiceUnavailable.WithDetails(enrollResp.Message)
	}

	// Step 4: Mark the request as completed
	if err := fprRepo.CompleteRequest(ctx, request.RequestID, request.CompanyID); err != nil {
		global.Logger.Error("Failed to complete request", err)
		// Don't return error since face was updated successfully
	}

	// Step 5: Invalidate token cache
	s.invalidateTokenCache(ctx, input.Token)

	return &appModel.UpdateFaceProfileOutput{
		Success:      true,
		ProfileID:    enrollResp.ProfileID,
		QualityScore: float64(enrollResp.QualityScore),
		Message:      "Face profile updated successfully",
	}, nil
}

// =================================
// Helper Methods:
// =================================

func (s *FaceProfileUpdateServiceImpl) checkAuthorization(session *appModel.SessionInfo, requestedCompanyID *string, minRole domainModel.Role, allowEmployeeSelfAccess bool) *appErrors.Error {
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

func (s *FaceProfileUpdateServiceImpl) checkSpamLocal(ctx context.Context, userID uuid.UUID) bool {
	localCache, err := domainCache.GetLocalCache()
	if err != nil {
		return false
	}

	key := constants.CacheKeyPrefixRequestLock + userID.String()
	exists, err := localCache.Exists(ctx, key)
	if err != nil {
		return false
	}
	return exists
}

func (s *FaceProfileUpdateServiceImpl) checkSpamDistributed(ctx context.Context, userID uuid.UUID) bool {
	distCache, err := domainCache.GetDistributedCache()
	if err != nil {
		return false
	}

	key := constants.CacheKeyPrefixRequestLock + userID.String()
	exists, err := distCache.Exists(ctx, key)
	if err != nil {
		return false
	}
	return exists
}

func (s *FaceProfileUpdateServiceImpl) setSpamMarker(ctx context.Context, userID uuid.UUID) {
	key := constants.CacheKeyPrefixRequestLock + userID.String()

	// Set in local cache
	if localCache, err := domainCache.GetLocalCache(); err == nil {
		_ = localCache.SetTTL(ctx, key, "1", int64(constants.SpamPreventionWindowSeconds))
	}

	// Set in distributed cache
	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.SetTTL(ctx, key, "1", int64(constants.SpamPreventionWindowSeconds))
	}
}

func (s *FaceProfileUpdateServiceImpl) checkPendingRequest(ctx context.Context, repo domainRepo.IFaceProfileUpdateRequestRepository, userID, companyID uuid.UUID) (bool, error) {
	// Try cache first
	key := constants.CacheKeyPrefixPendingRequest + userID.String()

	if localCache, err := domainCache.GetLocalCache(); err == nil {
		if val, err := localCache.Get(ctx, key); err == nil && val != "" {
			return val == "1", nil
		}
	}

	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		if val, err := distCache.Get(ctx, key); err == nil && val != "" {
			// Backfill local cache
			if localCache, err := domainCache.GetLocalCache(); err == nil {
				_ = localCache.SetTTL(ctx, key, val, int64(constants.TTLLocalPendingCheck))
			}
			return val == "1", nil
		}
	}

	// Query database
	hasPending, err := repo.HasPendingRequest(ctx, userID, companyID)
	if err != nil {
		return false, err
	}

	// Cache the result
	val := "0"
	if hasPending {
		val = "1"
	}

	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.SetTTL(ctx, key, val, int64(constants.TTLDistributedPendingCheck))
	}
	if localCache, err := domainCache.GetLocalCache(); err == nil {
		_ = localCache.SetTTL(ctx, key, val, int64(constants.TTLLocalPendingCheck))
	}

	return hasPending, nil
}

func (s *FaceProfileUpdateServiceImpl) invalidatePendingCache(ctx context.Context, userID, companyID uuid.UUID) {
	key := constants.CacheKeyPrefixPendingRequest + userID.String()

	if localCache, err := domainCache.GetLocalCache(); err == nil {
		_ = localCache.Delete(ctx, key)
	}
	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.Delete(ctx, key)
	}
}

func (s *FaceProfileUpdateServiceImpl) getMonthlyRequestCount(ctx context.Context, repo domainRepo.IFaceProfileUpdateRequestRepository, userID uuid.UUID, month string) (int, error) {
	key := constants.CacheKeyPrefixMonthlyCount + userID.String() + ":" + month

	// Try cache first
	if localCache, err := domainCache.GetLocalCache(); err == nil {
		if val, err := localCache.Get(ctx, key); err == nil && val != "" {
			var count int
			if _, parseErr := fmt.Sscanf(val, "%d", &count); parseErr == nil {
				return count, nil
			}
		}
	}

	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		if val, err := distCache.Get(ctx, key); err == nil && val != "" {
			var count int
			if _, parseErr := fmt.Sscanf(val, "%d", &count); parseErr == nil {
				// Backfill local cache
				if localCache, err := domainCache.GetLocalCache(); err == nil {
					_ = localCache.SetTTL(ctx, key, val, int64(constants.TTLLocalMonthlyCount))
				}
				return count, nil
			}
		}
	}

	// Query database
	count, err := repo.CountRequestsByUserInMonth(ctx, userID, month)
	if err != nil {
		return 0, err
	}

	// Cache the result
	countStr := fmt.Sprintf("%d", count)
	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.SetTTL(ctx, key, countStr, int64(constants.TTLDistributedMonthlyCount))
	}
	if localCache, err := domainCache.GetLocalCache(); err == nil {
		_ = localCache.SetTTL(ctx, key, countStr, int64(constants.TTLLocalMonthlyCount))
	}

	return count, nil
}

func (s *FaceProfileUpdateServiceImpl) invalidateMonthlyCountCache(ctx context.Context, userID uuid.UUID, month string) {
	key := constants.CacheKeyPrefixMonthlyCount + userID.String() + ":" + month

	if localCache, err := domainCache.GetLocalCache(); err == nil {
		_ = localCache.Delete(ctx, key)
	}
	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.Delete(ctx, key)
	}
}

func (s *FaceProfileUpdateServiceImpl) generateSecureToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return uuid.New().String()
	}
	return hex.EncodeToString(bytes)
}

func (s *FaceProfileUpdateServiceImpl) acquireApprovalLock(ctx context.Context, requestID uuid.UUID) bool {
	key := constants.CacheKeyPrefixApprovalLock + requestID.String()

	distCache, err := domainCache.GetDistributedCache()
	if err != nil {
		return true // Allow if cache is unavailable
	}

	// Use Lua script for atomic check-and-set
	script := `
		if redis.call("EXISTS", KEYS[1]) == 0 then
			redis.call("SET", KEYS[1], "1", "EX", 30)
			return 1
		end
		return 0
	`
	result, err := distCache.LuaScript(ctx, script, []string{key})
	if err != nil {
		return true // Allow if script fails
	}

	if val, ok := result.(int64); ok {
		return val == 1
	}
	return true
}

func (s *FaceProfileUpdateServiceImpl) releaseApprovalLock(ctx context.Context, requestID uuid.UUID) {
	key := constants.CacheKeyPrefixApprovalLock + requestID.String()

	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.Delete(ctx, key)
	}
}

func (s *FaceProfileUpdateServiceImpl) cacheUpdateToken(ctx context.Context, token, userID, companyID string, expiresAt time.Time) {
	key := constants.CacheKeyPrefixUpdateToken + token
	value := fmt.Sprintf("%s:%s:%d", userID, companyID, expiresAt.Unix())

	ttl := int64(expiresAt.Sub(time.Now()).Seconds())
	if ttl <= 0 {
		return
	}

	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.SetTTL(ctx, key, value, ttl)
	}
	if localCache, err := domainCache.GetLocalCache(); err == nil {
		localTTL := ttl
		if localTTL > int64(constants.TTLLocalUpdateToken) {
			localTTL = int64(constants.TTLLocalUpdateToken)
		}
		_ = localCache.SetTTL(ctx, key, value, localTTL)
	}
}

func (s *FaceProfileUpdateServiceImpl) getTokenFromCache(ctx context.Context, token string) (userID, companyID string, expiresAt time.Time, found bool) {
	key := constants.CacheKeyPrefixUpdateToken + token

	var value string

	// Try local cache first
	if localCache, err := domainCache.GetLocalCache(); err == nil {
		if val, err := localCache.Get(ctx, key); err == nil && val != "" {
			value = val
		}
	}

	// Try distributed cache
	if value == "" {
		if distCache, err := domainCache.GetDistributedCache(); err == nil {
			if val, err := distCache.Get(ctx, key); err == nil && val != "" {
				value = val
				// Backfill local cache
				if localCache, err := domainCache.GetLocalCache(); err == nil {
					_ = localCache.SetTTL(ctx, key, val, int64(constants.TTLLocalUpdateToken))
				}
			}
		}
	}

	if value == "" {
		return "", "", time.Time{}, false
	}

	var expiresUnix int64
	if _, err := fmt.Sscanf(value, "%s:%s:%d", &userID, &companyID, &expiresUnix); err != nil {
		// Try a different parsing approach using standard library
		parts := strings.SplitN(value, ":", 3)
		if len(parts) == 3 {
			userID = parts[0]
			companyID = parts[1]
			if _, err := fmt.Sscanf(parts[2], "%d", &expiresUnix); err != nil {
				return "", "", time.Time{}, false
			}
		} else {
			return "", "", time.Time{}, false
		}
	}

	return userID, companyID, time.Unix(expiresUnix, 0), true
}

func (s *FaceProfileUpdateServiceImpl) invalidateTokenCache(ctx context.Context, token string) {
	key := constants.CacheKeyPrefixUpdateToken + token

	if localCache, err := domainCache.GetLocalCache(); err == nil {
		_ = localCache.Delete(ctx, key)
	}
	if distCache, err := domainCache.GetDistributedCache(); err == nil {
		_ = distCache.Delete(ctx, key)
	}
}
