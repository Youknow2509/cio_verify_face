package handler

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/logger"
	pb "github.com/youknow2509/cio_verify_face/server/service_auth/proto"
)

// AuthGRPCHandler implements the gRPC AuthService
type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authCacheService service.IAuthCacheService
	coreAuthService  service.ICoreAuthService
	logger          logger.ILogger
}

// NewAuthGRPCHandler creates a new gRPC handler
func NewAuthGRPCHandler(
	authCacheService service.IAuthCacheService,
	coreAuthService service.ICoreAuthService,
	logger logger.ILogger,
) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		authCacheService: authCacheService,
		coreAuthService:  coreAuthService,
		logger:          logger,
	}
}

// ValidateToken validates an access token
func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &pb.ValidateTokenResponse{
			Valid:        false,
			ErrorMessage: "token is required",
		}, nil
	}

	// Use cache service for token validation
	result, err := h.authCacheService.ValidateAccessToken(ctx, req.Token)
	if err != nil {
		h.logger.Error("Failed to validate token", "error", err.Error())
		return &pb.ValidateTokenResponse{
			Valid:        false,
			ErrorMessage: err.Error(),
		}, nil
	}

	if !result.Valid {
		return &pb.ValidateTokenResponse{
			Valid:        false,
			ErrorMessage: "invalid token",
		}, nil
	}

	response := &pb.ValidateTokenResponse{
		Valid:     true,
		UserId:    result.UserID.String(),
		CompanyId: result.CompanyID.String(),
		SessionId: result.SessionID.String(),
		ExpiresAt: result.ExpiresAt.Unix(),
	}

	// Check company permission if required
	if req.CompanyId != nil {
		hasPermission, err := h.authCacheService.CheckUserPermissionCached(ctx, *req.CompanyId, result.UserID.String())
		if err != nil || !hasPermission {
			response.Valid = false
			response.ErrorMessage = "user does not have permission for this company"
			return response, nil
		}
	}

	// Check specific permission if required
	if req.RequiredPermission != nil {
		// This would require implementing permission checking logic
		// For now, we'll assume valid if token is valid and company check passed
		response.Permissions = result.Permissions
	}

	return response, nil
}

// GetUserInfo retrieves user information
func (h *AuthGRPCHandler) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	if req.UserId == "" {
		return &pb.GetUserInfoResponse{
			ErrorMessage: "user_id is required",
		}, nil
	}

	userInfo, err := h.authCacheService.GetUserInfoCached(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user info", "error", err.Error(), "user_id", req.UserId)
		return &pb.GetUserInfoResponse{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.GetUserInfoResponse{
		UserId:    userInfo.UserID,
		Email:     userInfo.Email,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		CompanyId: userInfo.CompanyID,
		Role:      userInfo.Role,
		IsActive:  userInfo.IsActive,
		CreatedAt: userInfo.CreatedAt.Unix(),
		UpdatedAt: userInfo.UpdatedAt.Unix(),
	}, nil
}

// CheckUserPermission checks if user has permission in a company
func (h *AuthGRPCHandler) CheckUserPermission(ctx context.Context, req *pb.CheckUserPermissionRequest) (*pb.CheckUserPermissionResponse, error) {
	if req.UserId == "" || req.CompanyId == "" {
		return &pb.CheckUserPermissionResponse{
			HasPermission: false,
			ErrorMessage:  "user_id and company_id are required",
		}, nil
	}

	hasPermission, err := h.authCacheService.CheckUserPermissionCached(ctx, req.CompanyId, req.UserId)
	if err != nil {
		h.logger.Error("Failed to check user permission", "error", err.Error(), "user_id", req.UserId, "company_id", req.CompanyId)
		return &pb.CheckUserPermissionResponse{
			HasPermission: false,
			ErrorMessage:  err.Error(),
		}, nil
	}

	// Get user info to include role information
	userInfo, err := h.authCacheService.GetUserInfoCached(ctx, req.UserId)
	response := &pb.CheckUserPermissionResponse{
		HasPermission: hasPermission,
	}

	if err == nil && userInfo != nil {
		response.UserRole = userInfo.Role
		// Add permissions based on role or from cache
		response.Permissions = []string{} // This would be populated based on your permission system
	}

	return response, nil
}

// CheckDeviceInCompany checks if a device belongs to a company
func (h *AuthGRPCHandler) CheckDeviceInCompany(ctx context.Context, req *pb.CheckDeviceInCompanyRequest) (*pb.CheckDeviceInCompanyResponse, error) {
	if req.DeviceId == "" || req.CompanyId == "" {
		return &pb.CheckDeviceInCompanyResponse{
			DeviceExists: false,
			ErrorMessage: "device_id and company_id are required",
		}, nil
	}

	exists, err := h.authCacheService.CheckDeviceInCompanyCached(ctx, req.CompanyId, req.DeviceId)
	if err != nil {
		h.logger.Error("Failed to check device in company", "error", err.Error(), "device_id", req.DeviceId, "company_id", req.CompanyId)
		return &pb.CheckDeviceInCompanyResponse{
			DeviceExists: false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.CheckDeviceInCompanyResponse{
		DeviceExists: exists,
		IsActive:     exists, // Assuming if device exists, it's active
		DeviceName:   "",     // This would require additional data retrieval
		DeviceType:   "",     // This would require additional data retrieval
	}, nil
}

// ValidateDeviceSession validates a device session token
func (h *AuthGRPCHandler) ValidateDeviceSession(ctx context.Context, req *pb.ValidateDeviceSessionRequest) (*pb.ValidateDeviceSessionResponse, error) {
	if req.DeviceId == "" || req.SessionToken == "" {
		return &pb.ValidateDeviceSessionResponse{
			Valid:        false,
			ErrorMessage: "device_id and session_token are required",
		}, nil
	}

	// Validate the session token - this would require implementing device session validation
	// For now, we'll use the regular token validation as a base
	result, err := h.authCacheService.ValidateAccessToken(ctx, req.SessionToken)
	if err != nil || !result.Valid {
		return &pb.ValidateDeviceSessionResponse{
			Valid:        false,
			ErrorMessage: "invalid session token",
		}, nil
	}

	// Additional device-specific validation would go here
	if req.CompanyId != "" {
		deviceExists, err := h.authCacheService.CheckDeviceInCompanyCached(ctx, req.CompanyId, req.DeviceId)
		if err != nil || !deviceExists {
			return &pb.ValidateDeviceSessionResponse{
				Valid:        false,
				ErrorMessage: "device not found in company",
			}, nil
		}
	}

	return &pb.ValidateDeviceSessionResponse{
		Valid:     true,
		DeviceId:  req.DeviceId,
		CompanyId: result.CompanyID.String(),
		UserId:    result.UserID.String(),
		ExpiresAt: result.ExpiresAt.Unix(),
	}, nil
}

// GetUserCompany retrieves company information for a user
func (h *AuthGRPCHandler) GetUserCompany(ctx context.Context, req *pb.GetUserCompanyRequest) (*pb.GetUserCompanyResponse, error) {
	if req.UserId == "" {
		return &pb.GetUserCompanyResponse{
			ErrorMessage: "user_id is required",
		}, nil
	}

	userInfo, err := h.authCacheService.GetUserInfoCached(ctx, req.UserId)
	if err != nil {
		h.logger.Error("Failed to get user company", "error", err.Error(), "user_id", req.UserId)
		return &pb.GetUserCompanyResponse{
			ErrorMessage: err.Error(),
		}, nil
	}

	// This would require additional company data retrieval
	// For now, we'll return basic information from user info
	return &pb.GetUserCompanyResponse{
		CompanyId:   userInfo.CompanyID,
		CompanyName: "", // This would require company service integration
		CompanyCode: "", // This would require company service integration
		IsActive:    userInfo.IsActive,
		UserRole:    userInfo.Role,
		Permissions: []string{}, // This would be populated based on role and permissions
	}, nil
}

// BatchValidateTokens validates multiple tokens in a single request
func (h *AuthGRPCHandler) BatchValidateTokens(ctx context.Context, req *pb.BatchValidateTokensRequest) (*pb.BatchValidateTokensResponse, error) {
	if len(req.Tokens) == 0 {
		return &pb.BatchValidateTokensResponse{
			Results: []*pb.TokenValidationResult{},
		}, nil
	}

	results := make([]*pb.TokenValidationResult, 0, len(req.Tokens))

	for _, tokenReq := range req.Tokens {
		result := &pb.TokenValidationResult{
			RequestId: tokenReq.RequestId,
			Valid:     false,
		}

		if tokenReq.Token == "" {
			result.ErrorMessage = "token is required"
			results = append(results, result)
			continue
		}

		// Validate token
		validationResult, err := h.authCacheService.ValidateAccessToken(ctx, tokenReq.Token)
		if err != nil {
			result.ErrorMessage = err.Error()
			results = append(results, result)
			continue
		}

		if !validationResult.Valid {
			result.ErrorMessage = "invalid token"
			results = append(results, result)
			continue
		}

		result.Valid = true
		result.UserId = validationResult.UserID.String()
		result.CompanyId = validationResult.CompanyID.String()
		result.SessionId = validationResult.SessionID.String()
		result.ExpiresAt = validationResult.ExpiresAt.Unix()
		result.Permissions = validationResult.Permissions

		// Check company permission if required
		if req.CompanyId != nil {
			hasPermission, err := h.authCacheService.CheckUserPermissionCached(ctx, *req.CompanyId, validationResult.UserID.String())
			if err != nil || !hasPermission {
				result.Valid = false
				result.ErrorMessage = "user does not have permission for this company"
			}
		}

		results = append(results, result)
	}

	return &pb.BatchValidateTokensResponse{
		Results: results,
	}, nil
}