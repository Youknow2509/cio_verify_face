package handler

import (
	"context"
	"errors"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/logger"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/uuid"
	pb "github.com/youknow2509/cio_verify_face/server/service_auth/proto"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// AuthGRPCHandler implements the gRPC AuthService
type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authCacheService service.IAuthCacheService
	coreAuthService  service.ICoreAuthService
	logger           logger.ILogger
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
		logger:           logger,
	}
}

func (a *AuthGRPCHandler) CreateUserToken(ctx context.Context, req *pb.CreateUserTokenRequest) (*pb.CreateUserTokenResponse, error) {
	tok := service.GetTokenService()
	if tok == nil {
		return nil, status.Error(codes.FailedPrecondition, "token service not initialized")
	}

	in, err := toModelCreateUserTokenInput(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	out, err := tok.CreateUserToken(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create user token failed: %v", err)
	}

	resp, err := toPbCreateUserTokenResponse(out)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response mapping failed: %v", err)
	}
	return resp, nil
}

func (a *AuthGRPCHandler) CreateServiceToken(ctx context.Context, req *pb.CreateServiceTokenRequest) (*pb.CreateServiceTokenResponse, error) {
	// TokenService hiện không hỗ trợ service token theo interface hiện có.
	return nil, status.Errorf(codes.Unimplemented, "service token is not supported")
}

func (a *AuthGRPCHandler) CreateDeviceToken(ctx context.Context, req *pb.CreateDeviceTokenRequest) (*pb.CreateDeviceTokenResponse, error) {
	tok := service.GetTokenService()
	if tok == nil {
		return nil, status.Error(codes.FailedPrecondition, "token service not initialized")
	}

	in, err := toModelCreateDeviceTokenInput(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	token, err := tok.CreateTokenDevice(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create device token failed: %v", err)
	}

	resp, err := toPbCreateDeviceTokenResponse(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response mapping failed: %v", err)
	}
	return resp, nil
}

func (a *AuthGRPCHandler) ParseUserToken(ctx context.Context, req *pb.ParseUserTokenRequest) (*pb.ParseUserTokenResponse, error) {
	tok := service.GetTokenService()
	if tok == nil {
		return nil, status.Error(codes.FailedPrecondition, "token service not initialized")
	}

	in, err := toModelParseUserTokenInput(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	out, err := tok.ParseTokenUser(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "parse user token failed: %v", err)
	}

	resp, err := toPbParseUserTokenResponse(out)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response mapping failed: %v", err)
	}
	return resp, nil
}

func (a *AuthGRPCHandler) ParseServiceToken(ctx context.Context, req *pb.ParseServiceTokenRequest) (*pb.ParseServiceTokenResponse, error) {
	// TokenService hiện không hỗ trợ service token theo interface hiện có.
	return nil, status.Errorf(codes.Unimplemented, "service token is not supported")
}

func (a *AuthGRPCHandler) ParseDeviceToken(ctx context.Context, req *pb.ParseDeviceTokenRequest) (*pb.ParseDeviceTokenResponse, error) {
	tok := service.GetTokenService()
	if tok == nil {
		return nil, status.Error(codes.FailedPrecondition, "token service not initialized")
	}

	in, err := toModelCheckDeviceTokenInput(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	valid, deviceID, err := tok.CheckTokenDevice(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "parse device token failed: %v", err)
	}

	resp, err := toPbParseDeviceTokenResponse(valid, deviceID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response mapping failed: %v", err)
	}
	return resp, nil
}

// -------- Mapping helpers (TODO: map field theo proto/model thực tế) --------
func toModelCreateUserTokenInput(req *pb.CreateUserTokenRequest) (model.CreateTokenUserInput, error) {
	userUuid, err := uuid.ParseUUID(req.GetUserId())
	if err != nil {
		return model.CreateTokenUserInput{}, errors.New("invalid user ID format")
	}
	return model.CreateTokenUserInput{
		UserId: userUuid,
		Role:   int(req.GetRoles()),
	}, nil
}

func toPbCreateUserTokenResponse(out *model.CreateTokenUserOutput) (*pb.CreateUserTokenResponse, error) {
	return &pb.CreateUserTokenResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, nil
}

func toModelCreateDeviceTokenInput(req *pb.CreateDeviceTokenRequest) (model.CreateTokenDeviceInput, error) {
	deviceUuid, err := uuid.ParseUUID(req.GetDeviceId())
	if err != nil {
		return model.CreateTokenDeviceInput{}, errors.New("invalid device ID format")
	}
	companyUuid, err := uuid.ParseUUID(req.GetCompanyId())
	if err != nil {
		return model.CreateTokenDeviceInput{}, errors.New("invalid company ID format")
	}
	return model.CreateTokenDeviceInput{
		DeviceId:  deviceUuid,
		CompanyId: companyUuid,
	}, nil
}

func toPbCreateDeviceTokenResponse(token string) (*pb.CreateDeviceTokenResponse, error) {
	return &pb.CreateDeviceTokenResponse{
		Token: token,
	}, nil
}

func toModelParseUserTokenInput(req *pb.ParseUserTokenRequest) (model.ParseTokenUserInput, error) {
	return model.ParseTokenUserInput{
		Token: req.GetToken(),
	}, nil
}

func toPbParseUserTokenResponse(out *model.ParseTokenUserOutput) (*pb.ParseUserTokenResponse, error) {
	return &pb.ParseUserTokenResponse{
		UserId: out.UserId,
		Roles:  int32(out.Role),
	}, nil
}

func toModelCheckDeviceTokenInput(req *pb.ParseDeviceTokenRequest) (model.CheckTokenDeviceInput, error) {
	deviceUuid, err := uuid.ParseUUID(req.GetDeviceId())
	if err != nil {
		return model.CheckTokenDeviceInput{}, errors.New("invalid device ID format")
	}
	return model.CheckTokenDeviceInput{
		Token:    req.GetToken(),
		DeviceId: deviceUuid,
	}, nil
}

func toPbParseDeviceTokenResponse(valid bool, deviceID string) (*pb.ParseDeviceTokenResponse, error) {
	return &pb.ParseDeviceTokenResponse{
		DeviceId: deviceID,
	}, nil
}
