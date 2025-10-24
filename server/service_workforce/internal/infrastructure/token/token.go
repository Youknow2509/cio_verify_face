package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	domainErrors "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/errors"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/model"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/token"
	pb "github.com/youknow2509/cio_verify_face/server/service_workforce/proto"
)

// =======================================================
// Define token infrastructure implementation - Use grpc call to token service
// =======================================================
type TokenService struct {
	grpc pb.AuthServiceClient
}

// CreateUserToken implements token.ITokenService.
func (t *TokenService) CreateUserToken(ctx context.Context, input *domainModel.TokenUserJwtInput) (*domainModel.UserTokenOutput, error) {
	tokenResp, err := t.grpc.CreateUserToken(
		ctx,
		&pb.CreateUserTokenRequest{
			UserId: input.UserId,
			Roles:  int32(input.Role),
		},
	)
	if err != nil {
		return nil, err
	}
	return &domainModel.UserTokenOutput{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}, nil
}

// ParseUserRefreshToken implements token.ITokenService.
func (t *TokenService) ParseUserRefreshToken(ctx context.Context, token string) (*domainModel.TokenUserRefreshOutput, *domainErrors.TokenValidationError) {
	panic("unimplemented")
}

func (t *TokenService) ParseUserToken(ctx context.Context, token string) (*domainModel.TokenUserJwtOutput, *domainErrors.TokenValidationError) {
	resp, err := t.grpc.ParseUserToken(ctx, &pb.ParseUserTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, handleError(err)
	}
	return &domainModel.TokenUserJwtOutput{
		UserId:    resp.UserId,
		Role:      int(resp.Roles),
		TokenId:   resp.TokenId,
		ExpiresAt: time.Unix(resp.ExpriresAt, 0),
	}, nil
}

// CheckDeviceToken implements token.ITokenService.
func (t *TokenService) CheckDeviceToken(ctx context.Context, token string) (bool, *domainErrors.TokenValidationError) {
	_, err := t.grpc.ParseDeviceToken(ctx, &pb.ParseDeviceTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, handleError(err)
	}
	return true, nil
}

// CreateDeviceToken implements token.ITokenService.
func (t *TokenService) CreateDeviceToken(ctx context.Context, input *domainModel.TokenDeviceJwtInput) (string, error) {
	token, err := t.grpc.CreateDeviceToken(ctx, &pb.CreateDeviceTokenRequest{
		DeviceId:  input.DeviceId,
		CompanyId: input.CompanyId,
	})
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

/**
 * New token implementation
 */
func NewTokenService(
	grpcClient pb.AuthServiceClient,
) domainToken.ITokenService {
	return &TokenService{
		grpc: grpcClient,
	}
}

// =======================================================
//
//	Helper functions
//
// =======================================================
func handleError(err error) *domainErrors.TokenValidationError {
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return domainErrors.GetTokenValidationError(domainErrors.TokenMalformedErrorCode)
	}
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return domainErrors.GetTokenValidationError(domainErrors.TokenSignatureInvalidErrCode)
	}
	if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return domainErrors.GetTokenValidationError(domainErrors.TokenExpiredErrorCode)
	}
	if err != nil {
		return domainErrors.NewTokenServiceValidationError(err)
	}
	return nil
}
