package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/errors"
	appModel "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/service/impl"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/constants"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/model"
)

// =================================
// Mock Implementations:
// =================================

// MockFaceProfileUpdateRequestRepository
type MockFaceProfileUpdateRequestRepository struct {
	requests       map[string]*domainModel.FaceProfileUpdateRequest
	pendingCheck   map[string]bool
	monthlyCount   map[string]int
}

func NewMockFaceProfileUpdateRequestRepository() *MockFaceProfileUpdateRequestRepository {
	return &MockFaceProfileUpdateRequestRepository{
		requests:     make(map[string]*domainModel.FaceProfileUpdateRequest),
		pendingCheck: make(map[string]bool),
		monthlyCount: make(map[string]int),
	}
}

func (m *MockFaceProfileUpdateRequestRepository) CreateRequest(ctx context.Context, req *domainModel.FaceProfileUpdateRequest) error {
	key := req.RequestID.String()
	m.requests[key] = req
	m.pendingCheck[req.UserID.String()] = true
	monthKey := req.UserID.String() + ":" + req.RequestMonth
	m.monthlyCount[monthKey]++
	return nil
}

func (m *MockFaceProfileUpdateRequestRepository) GetRequestByID(ctx context.Context, requestID, companyID uuid.UUID) (*domainModel.FaceProfileUpdateRequest, error) {
	return m.requests[requestID.String()], nil
}

func (m *MockFaceProfileUpdateRequestRepository) GetRequestByToken(ctx context.Context, token string) (*domainModel.FaceProfileUpdateRequest, error) {
	for _, req := range m.requests {
		if req.UpdateToken != nil && *req.UpdateToken == token {
			return req, nil
		}
	}
	return nil, nil
}

func (m *MockFaceProfileUpdateRequestRepository) GetPendingRequestsByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domainModel.FaceProfileUpdateRequest, error) {
	var result []*domainModel.FaceProfileUpdateRequest
	for _, req := range m.requests {
		if req.CompanyID == companyID && req.Status == domainModel.RequestStatusPending {
			result = append(result, req)
		}
	}
	return result, nil
}

func (m *MockFaceProfileUpdateRequestRepository) GetRequestsByUserAndMonth(ctx context.Context, userID uuid.UUID, month string) ([]*domainModel.FaceProfileUpdateRequest, error) {
	var result []*domainModel.FaceProfileUpdateRequest
	for _, req := range m.requests {
		if req.UserID == userID && req.RequestMonth == month {
			result = append(result, req)
		}
	}
	return result, nil
}

func (m *MockFaceProfileUpdateRequestRepository) CountRequestsByUserInMonth(ctx context.Context, userID uuid.UUID, month string) (int, error) {
	key := userID.String() + ":" + month
	return m.monthlyCount[key], nil
}

func (m *MockFaceProfileUpdateRequestRepository) UpdateRequestStatus(ctx context.Context, requestID, companyID uuid.UUID, status domainModel.RequestStatus) error {
	if req, ok := m.requests[requestID.String()]; ok {
		req.Status = status
		req.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockFaceProfileUpdateRequestRepository) ApproveRequest(ctx context.Context, requestID, companyID, approvedBy uuid.UUID, updateToken string, expiresAt interface{}) error {
	if req, ok := m.requests[requestID.String()]; ok {
		req.Status = domainModel.RequestStatusApproved
		req.UpdateToken = &updateToken
		req.ApprovedBy = &approvedBy
		now := time.Now()
		req.ApprovedAt = &now
		if exp, ok := expiresAt.(time.Time); ok {
			req.UpdateLinkExpiresAt = &exp
		}
		req.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockFaceProfileUpdateRequestRepository) RejectRequest(ctx context.Context, requestID, companyID, rejectedBy uuid.UUID, reason string) error {
	if req, ok := m.requests[requestID.String()]; ok {
		req.Status = domainModel.RequestStatusRejected
		req.RejectionReason = &reason
		req.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockFaceProfileUpdateRequestRepository) CompleteRequest(ctx context.Context, requestID, companyID uuid.UUID) error {
	if req, ok := m.requests[requestID.String()]; ok {
		req.Status = domainModel.RequestStatusCompleted
		req.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockFaceProfileUpdateRequestRepository) MarkExpiredRequests(ctx context.Context) (int64, error) {
	var count int64
	for _, req := range m.requests {
		if req.Status == domainModel.RequestStatusApproved && req.UpdateLinkExpiresAt != nil && time.Now().After(*req.UpdateLinkExpiresAt) {
			req.Status = domainModel.RequestStatusExpired
			count++
		}
	}
	return count, nil
}

func (m *MockFaceProfileUpdateRequestRepository) HasPendingRequest(ctx context.Context, userID, companyID uuid.UUID) (bool, error) {
	for _, req := range m.requests {
		if req.UserID == userID && req.CompanyID == companyID && req.Status == domainModel.RequestStatusPending {
			return true, nil
		}
	}
	return false, nil
}

// =================================
// Unit Tests:
// =================================

func TestRequestStatus(t *testing.T) {
	// Test status constants
	tests := []struct {
		status   domainModel.RequestStatus
		expected int
	}{
		{domainModel.RequestStatusPending, 0},
		{domainModel.RequestStatusApproved, 1},
		{domainModel.RequestStatusRejected, 2},
		{domainModel.RequestStatusExpired, 3},
		{domainModel.RequestStatusCompleted, 4},
	}

	for _, test := range tests {
		if int(test.status) != test.expected {
			t.Errorf("Expected status %d but got %d", test.expected, int(test.status))
		}
	}
}

func TestRoleConstants(t *testing.T) {
	// Test role constants match expected values
	tests := []struct {
		role     domainModel.Role
		expected int
	}{
		{domainModel.RoleSystemAdmin, 0},
		{domainModel.RoleCompanyAdmin, 1},
		{domainModel.RoleEmployee, 2},
	}

	for _, test := range tests {
		if int(test.role) != test.expected {
			t.Errorf("Expected role %d but got %d", test.expected, int(test.role))
		}
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have reasonable values
	if constants.MaxFaceProfileRequestsPerMonth <= 0 {
		t.Error("MaxFaceProfileRequestsPerMonth should be positive")
	}

	if constants.UpdateLinkTTLSeconds <= 0 {
		t.Error("UpdateLinkTTLSeconds should be positive")
	}

	if constants.PasswordResetCooldownSeconds <= 0 {
		t.Error("PasswordResetCooldownSeconds should be positive")
	}

	if constants.MaxPasswordResetsPerManagerPerHour <= 0 {
		t.Error("MaxPasswordResetsPerManagerPerHour should be positive")
	}
}

func TestErrorCodes(t *testing.T) {
	// Test that error codes are unique
	codes := map[int]string{}
	errors := []*appErrors.Error{
		appErrors.ErrUnknown,
		appErrors.ErrInvalidInput,
		appErrors.ErrUnauthorized,
		appErrors.ErrForbidden,
		appErrors.ErrNotFound,
		appErrors.ErrConflict,
		appErrors.ErrServiceUnavailable,
		appErrors.ErrRateLimitExceeded,
		appErrors.ErrSpamDetected,
		appErrors.ErrMonthlyLimitReached,
		appErrors.ErrRequestAlreadyPending,
		appErrors.ErrRequestNotFound,
		appErrors.ErrRequestAlreadyProcessed,
		appErrors.ErrInvalidUpdateToken,
		appErrors.ErrUpdateTokenExpired,
		appErrors.ErrFaceEnrollmentFailed,
		appErrors.ErrPasswordResetSpam,
		appErrors.ErrEmployeeNotFound,
		appErrors.ErrPasswordResetFailed,
		appErrors.ErrKafkaPublishFailed,
	}

	for _, err := range errors {
		if existing, ok := codes[err.Code]; ok {
			t.Errorf("Duplicate error code %d: %s and %s", err.Code, existing, err.Message)
		}
		codes[err.Code] = err.Message
	}
}

func TestErrorWithDetails(t *testing.T) {
	err := appErrors.ErrInvalidInput.WithDetails("test details")
	
	if err.Details != "test details" {
		t.Errorf("Expected details 'test details' but got '%s'", err.Details)
	}
	
	if err.Code != appErrors.ErrCodeInvalidInput {
		t.Errorf("Expected code %d but got %d", appErrors.ErrCodeInvalidInput, err.Code)
	}
	
	// Original error should not be modified
	if appErrors.ErrInvalidInput.Details != "" {
		t.Error("Original error should not be modified")
	}
}

func TestSessionInfoToUserInfo(t *testing.T) {
	userID := uuid.New()
	companyID := uuid.New()
	
	session := &appModel.SessionInfo{
		UserID:      userID.String(),
		CompanyID:   companyID.String(),
		Role:        int32(domainModel.RoleCompanyAdmin),
		SessionID:   uuid.New().String(),
		ClientIP:    "127.0.0.1",
		ClientAgent: "TestAgent",
	}
	
	userInfo := session.ToUserInfo()
	
	if userInfo.UserID != userID {
		t.Errorf("Expected user ID %s but got %s", userID, userInfo.UserID)
	}
	
	if *userInfo.CompanyID != companyID {
		t.Errorf("Expected company ID %s but got %s", companyID, *userInfo.CompanyID)
	}
	
	if userInfo.Role != domainModel.RoleCompanyAdmin {
		t.Errorf("Expected role %d but got %d", domainModel.RoleCompanyAdmin, userInfo.Role)
	}
}

func TestFaceProfileUpdateRequestDTO(t *testing.T) {
	requestID := uuid.New()
	userID := uuid.New()
	companyID := uuid.New()
	now := time.Now()
	token := "test-token"
	expiresAt := now.Add(24 * time.Hour)
	
	req := &domainModel.FaceProfileUpdateRequest{
		RequestID:           requestID,
		UserID:              userID,
		CompanyID:           companyID,
		Status:              domainModel.RequestStatusApproved,
		RequestMonth:        "2024-01",
		RequestCountInMonth: 1,
		UpdateToken:         &token,
		UpdateLinkExpiresAt: &expiresAt,
		CreatedAt:           now,
		UpdatedAt:           now,
		MetaData:            map[string]interface{}{},
	}
	
	baseURL := "https://example.com"
	dto := appModel.ToFaceProfileUpdateRequestDTO(req, baseURL)
	
	if dto.RequestID != requestID.String() {
		t.Errorf("Expected request ID %s but got %s", requestID.String(), dto.RequestID)
	}
	
	if dto.Status != "approved" {
		t.Errorf("Expected status 'approved' but got '%s'", dto.Status)
	}
	
	expectedLink := baseURL + "/api/v1/profile-update/face?token=" + token
	if dto.UpdateLink != expectedLink {
		t.Errorf("Expected update link '%s' but got '%s'", expectedLink, dto.UpdateLink)
	}
}

func TestMockRepository(t *testing.T) {
	repo := NewMockFaceProfileUpdateRequestRepository()
	ctx := context.Background()
	
	userID := uuid.New()
	companyID := uuid.New()
	requestID := uuid.New()
	
	// Test CreateRequest
	req := &domainModel.FaceProfileUpdateRequest{
		RequestID:    requestID,
		UserID:       userID,
		CompanyID:    companyID,
		Status:       domainModel.RequestStatusPending,
		RequestMonth: "2024-01",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	err := repo.CreateRequest(ctx, req)
	if err != nil {
		t.Errorf("CreateRequest failed: %v", err)
	}
	
	// Test GetRequestByID
	retrieved, err := repo.GetRequestByID(ctx, requestID, companyID)
	if err != nil {
		t.Errorf("GetRequestByID failed: %v", err)
	}
	if retrieved == nil {
		t.Error("GetRequestByID returned nil")
	}
	if retrieved.RequestID != requestID {
		t.Errorf("Expected request ID %s but got %s", requestID, retrieved.RequestID)
	}
	
	// Test HasPendingRequest
	hasPending, err := repo.HasPendingRequest(ctx, userID, companyID)
	if err != nil {
		t.Errorf("HasPendingRequest failed: %v", err)
	}
	if !hasPending {
		t.Error("Expected HasPendingRequest to return true")
	}
	
	// Test CountRequestsByUserInMonth
	count, err := repo.CountRequestsByUserInMonth(ctx, userID, "2024-01")
	if err != nil {
		t.Errorf("CountRequestsByUserInMonth failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1 but got %d", count)
	}
	
	// Test ApproveRequest
	token := "test-approval-token"
	expiresAt := time.Now().Add(24 * time.Hour)
	approverID := uuid.New()
	
	err = repo.ApproveRequest(ctx, requestID, companyID, approverID, token, expiresAt)
	if err != nil {
		t.Errorf("ApproveRequest failed: %v", err)
	}
	
	retrieved, _ = repo.GetRequestByID(ctx, requestID, companyID)
	if retrieved.Status != domainModel.RequestStatusApproved {
		t.Errorf("Expected status approved but got %d", retrieved.Status)
	}
	if retrieved.UpdateToken == nil || *retrieved.UpdateToken != token {
		t.Error("Update token not set correctly")
	}
	
	// Test GetRequestByToken
	retrievedByToken, err := repo.GetRequestByToken(ctx, token)
	if err != nil {
		t.Errorf("GetRequestByToken failed: %v", err)
	}
	if retrievedByToken == nil {
		t.Error("GetRequestByToken returned nil")
	}
	if retrievedByToken.RequestID != requestID {
		t.Errorf("Expected request ID %s but got %s", requestID, retrievedByToken.RequestID)
	}
}

func TestFaceProfileUpdateServiceImpl(t *testing.T) {
	// Test that service can be created
	service := impl.NewFaceProfileUpdateService("https://example.com")
	if service == nil {
		t.Error("NewFaceProfileUpdateService returned nil")
	}
}

func TestPasswordResetServiceImpl(t *testing.T) {
	// Test that service can be created
	service := impl.NewPasswordResetService()
	if service == nil {
		t.Error("NewPasswordResetService returned nil")
	}
}

// TestCreateUpdateRequestWithoutSession tests that creating a request without session fails
func TestCreateUpdateRequestWithoutSession(t *testing.T) {
	service := impl.NewFaceProfileUpdateService("https://example.com")
	ctx := context.Background()
	
	input := &appModel.CreateFaceProfileUpdateRequestInput{
		Session: nil,
		Reason:  "Test reason",
	}
	
	_, err := service.CreateUpdateRequest(ctx, input)
	if err == nil {
		t.Error("Expected error when session is nil")
	}
	if err.Code != appErrors.ErrCodeUnauthorized {
		t.Errorf("Expected unauthorized error code but got %d", err.Code)
	}
}

// TestGetMyUpdateRequestsWithoutSession tests that getting requests without session fails
func TestGetMyUpdateRequestsWithoutSession(t *testing.T) {
	service := impl.NewFaceProfileUpdateService("https://example.com")
	ctx := context.Background()
	
	input := &appModel.GetMyUpdateRequestsInput{
		Session: nil,
	}
	
	_, err := service.GetMyUpdateRequests(ctx, input)
	if err == nil {
		t.Error("Expected error when session is nil")
	}
	if err.Code != appErrors.ErrCodeUnauthorized {
		t.Errorf("Expected unauthorized error code but got %d", err.Code)
	}
}

// TestValidateUpdateTokenEmpty tests that validating empty token fails
func TestValidateUpdateTokenEmpty(t *testing.T) {
	service := impl.NewFaceProfileUpdateService("https://example.com")
	ctx := context.Background()
	
	input := &appModel.ValidateUpdateTokenInput{
		Token: "",
	}
	
	_, err := service.ValidateUpdateToken(ctx, input)
	if err == nil {
		t.Error("Expected error when token is empty")
	}
	if err.Code != appErrors.ErrCodeInvalidInput {
		t.Errorf("Expected invalid input error code but got %d", err.Code)
	}
}

// TestPasswordResetWithoutSession tests that password reset without session fails
func TestPasswordResetWithoutSession(t *testing.T) {
	service := impl.NewPasswordResetService()
	ctx := context.Background()
	
	input := &appModel.ResetEmployeePasswordInput{
		Session:    nil,
		EmployeeID: uuid.New().String(),
	}
	
	_, err := service.ResetEmployeePassword(ctx, input)
	if err == nil {
		t.Error("Expected error when session is nil")
	}
	if err.Code != appErrors.ErrCodeUnauthorized {
		t.Errorf("Expected unauthorized error code but got %d", err.Code)
	}
}
