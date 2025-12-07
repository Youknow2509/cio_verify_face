package handler

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/logger"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
	uuidUtils "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/uuid"
	pb "github.com/youknow2509/cio_verify_face/server/service_auth/proto"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (a *AuthGRPCHandler) HealthCheck(ctx context.Context, rep *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
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
	global.Logger.Info("Parse access token: ", req.GetToken())
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

	output, err := tok.ParseTokenDevice(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "parse device token failed: %v", err)
	}
	if output == nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid device token")
	}

	resp, err := toPbParseDeviceTokenResponse(output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response mapping failed: %v", err)
	}
	return resp, nil
}

// -------- Mapping helpers (TODO: map field theo proto/model thực tế) --------
func toModelCreateUserTokenInput(req *pb.CreateUserTokenRequest) (model.CreateTokenUserInput, error) {
	userUuid, err := uuidUtils.ParseUUID(req.GetUserId())
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
	deviceUuid, err := uuidUtils.ParseUUID(req.GetDeviceId())
	if err != nil {
		return model.CreateTokenDeviceInput{}, errors.New("invalid device ID format")
	}
	return model.CreateTokenDeviceInput{
		DeviceId:  deviceUuid,
		CompanyId: uuid.Nil,
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
		UserId:     out.UserId,
		Roles:      int32(out.Role),
		CompanyId:  out.CompanyId,
		TokenId:    out.TokenId,
		ExpriresAt: out.Expires.Unix(),
	}, nil
}

func toModelCheckDeviceTokenInput(req *pb.ParseDeviceTokenRequest) (model.ParseTokenDeviceInput, error) {
	return model.ParseTokenDeviceInput{
		Token: req.GetToken(),
	}, nil
}

func toPbParseDeviceTokenResponse(input *model.ParseTokenDeviceOutput) (*pb.ParseDeviceTokenResponse, error) {
	return &pb.ParseDeviceTokenResponse{
		DeviceId:  input.DeviceId,
		TokenId:   input.TokenId,
		CompanyId: input.CompanyId,
		ExpiresAt: input.Expires.Unix(),
	}, nil
}
