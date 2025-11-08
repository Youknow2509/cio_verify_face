package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	domainErrors "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/errors"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/token"
	utilsRandom "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/random"
)

// =======================================================
//
//	Token claims
//
// =======================================================
type (
	TokenUserJwtClaim struct {
		jwt.RegisteredClaims
		UserId string `json:"user_id"`
		Role   int    `json:"role"`
	}

	TokenUserRefreshJwtClaim struct {
		jwt.RegisteredClaims
	}

	TokenServiceJwtClaim struct {
		jwt.RegisteredClaims
		ServiceId   string `json:"service_id"`
		ServiceName string `json:"service_name"`
		Type        int    `json:"type"`
	}

	TokenDeviceJwtClaim struct {
		jwt.RegisteredClaims
		DeviceId  string `json:"device_id"`
		CompanyId string `json:"company_id"`
	}
)

// =======================================================
// Define token infrastructure implementation in domain
// =======================================================
type TokenService struct {
	secret   string
	issuer   string
	subject  string
	audience []string
}

// ParseUserRefreshToken implements token.ITokenService.
func (t *TokenService) ParseUserRefreshToken(ctx context.Context, token string) (*domainModel.TokenUserRefreshOutput, *domainErrors.TokenValidationError) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&TokenUserRefreshJwtClaim{},
		func(token *jwt.Token) (any, error) {
			return []byte(t.secret), nil
		},
	)
	if e := handleError(err); e != nil {
		return nil, e
	}
	if !parsedToken.Valid {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenValidationErrorCode)
	}
	out := parsedToken.Claims.(*TokenUserRefreshJwtClaim)
	if out == nil {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenMalformedErrorCode)
	}
	output := &domainModel.TokenUserRefreshOutput{
		TokenId:   out.ID,
		Issuer:    out.Issuer,
		Subject:   out.Subject,
		Audience:  out.Audience,
		ExpiresAt: out.ExpiresAt.Time,
		IssuedAt:  out.IssuedAt.Time,
	}
	return output, nil
}

// CreateDeviceRefreshToken implements token.ITokenService.
func (t *TokenService) CreateDeviceRefreshToken(ctx context.Context, input *domainModel.TokenDeviceRefreshInput) (string, error) {
	token := fmt.Sprintf(
		"%s.%s.%s.%s",
		utilsRandom.RandomString(8),
		utilsRandom.RandomString(8),
		utilsRandom.RandomString(8),
		utilsRandom.RandomString(8),
	)
	return token, nil
}

// CreateDeviceToken implements token.ITokenService.
func (t *TokenService) CreateDeviceToken(ctx context.Context, input *domainModel.TokenDeviceJwtInput) (string, error) {
	tokenClaims := &TokenDeviceJwtClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   t.subject,
			Audience:  t.audience,
			ExpiresAt: jwt.NewNumericDate(input.Expires),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        input.TokenId,
		},
		DeviceId:  input.DeviceId,
		CompanyId: input.CompanyId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString([]byte(t.secret))
}

// CreateServiceRefreshToken implements token.ITokenService.
func (t *TokenService) CreateServiceRefreshToken(ctx context.Context, input *domainModel.ServiceRefreshTokenInput) (string, error) {
	token := fmt.Sprintf(
		"%s.%s.%s.%s",
		utilsRandom.RandomString(8),
		utilsRandom.RandomString(8),
		utilsRandom.RandomString(8),
		utilsRandom.RandomString(8),
	)
	return token, nil
}

// ParseDeviceToken implements token.ITokenService.
func (t *TokenService) ParseDeviceToken(ctx context.Context, token string) (*domainModel.TokenDeviceJwtOutput, *domainErrors.TokenValidationError) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&TokenDeviceJwtClaim{},
		func(token *jwt.Token) (any, error) {
			return []byte(t.secret), nil
		},
	)
	if e := handleError(err); e != nil {
		return nil, e
	}
	if !parsedToken.Valid {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenValidationErrorCode)
	}
	out := parsedToken.Claims.(*TokenDeviceJwtClaim)
	if out == nil {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenMalformedErrorCode)
	}
	output := &domainModel.TokenDeviceJwtOutput{
		DeviceId:  out.DeviceId,
		CompanyId: out.CompanyId,
		TokenId:   out.ID,
		Issuer:    out.Issuer,
		Subject:   out.Subject,
		Audience:  out.Audience,
		ExpiresAt: out.ExpiresAt.Time,
		IssuedAt:  out.IssuedAt.Time,
		NotBefore: out.NotBefore.Time,
	}
	return output, nil
}

// CreateServiceToken implements token.ITokenService.
func (t *TokenService) CreateServiceToken(ctx context.Context, input *domainModel.TokenServiceJwtInput) (string, error) {
	tokenClaims := &TokenServiceJwtClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   t.subject,
			Audience:  t.audience,
			ExpiresAt: jwt.NewNumericDate(input.Expires),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        input.TokenId,
		},
		ServiceName: input.ServiceName,
		Type:        input.Type,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString([]byte(t.secret))
}

// CreateUserToken implements token.ITokenService.
func (t *TokenService) CreateUserToken(ctx context.Context, input *domainModel.TokenUserJwtInput) (string, error) {
	tokenClaim := &TokenUserJwtClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   t.subject,
			Audience:  t.audience,
			ExpiresAt: jwt.NewNumericDate(input.Expires),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        input.TokenId,
		},
		UserId: input.UserId,
		Role:   input.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	return token.SignedString([]byte(t.secret))
}

// CreateUserRefreshToken implements token.ITokenService.
func (t *TokenService) CreateUserRefreshToken(ctx context.Context, input *domainModel.TokenUserRefreshInput) (string, error) {
	tokenClaim := &TokenUserRefreshJwtClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   t.subject,
			Audience:  t.audience,
			ExpiresAt: jwt.NewNumericDate(input.Expires),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        input.TokenId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	return token.SignedString([]byte(t.secret))
}

// ParseServiceToken implements token.ITokenService.
func (t *TokenService) ParseServiceToken(ctx context.Context, token string) (*domainModel.TokenServiceJwtOutput, *domainErrors.TokenValidationError) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&TokenServiceJwtClaim{},
		func(token *jwt.Token) (any, error) {
			return []byte(t.secret), nil
		},
	)
	if e := handleError(err); e != nil {
		return nil, e
	}
	if !parsedToken.Valid {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenValidationErrorCode)
	}
	out := parsedToken.Claims.(*TokenServiceJwtClaim)
	if out == nil {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenMalformedErrorCode)
	}
	output := &domainModel.TokenServiceJwtOutput{
		ServiceName: out.ServiceName,
		Type:        out.Type,
		TokenId:     out.ID,
		Issuer:      out.Issuer,
		Subject:     out.Subject,
		Audience:    out.Audience,
		ExpiresAt:   out.ExpiresAt.Time,
		IssuedAt:    out.IssuedAt.Time,
		ServiceId:   out.ServiceId,
	}
	return output, nil
}

// ParseUserToken implements token.ITokenService.
func (t *TokenService) ParseUserToken(ctx context.Context, tokenString string) (*domainModel.TokenUserJwtOutput, *domainErrors.TokenValidationError) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenUserJwtClaim{},
		func(token *jwt.Token) (any, error) {
			return []byte(t.secret), nil
		},
	)
	if e := handleError(err); e != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Return token even if expired
			out := token.Claims.(*TokenUserJwtClaim)
			if out == nil {
				return nil, domainErrors.GetTokenValidationError(domainErrors.TokenMalformedErrorCode)
			}
			output := &domainModel.TokenUserJwtOutput{
				UserId:    out.UserId,
				TokenId:   out.ID,
				Role:      out.Role,
				Issuer:    out.Issuer,
				Subject:   out.Subject,
				Audience:  out.Audience,
				ExpiresAt: out.ExpiresAt.Time,
				IssuedAt:  out.IssuedAt.Time,
				NotBefore: out.NotBefore.Time,
			}
			return output, nil
		} else {
			return nil, e
		}
	}
	if !token.Valid {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenValidationErrorCode)
	}
	out := token.Claims.(*TokenUserJwtClaim)
	if out == nil {
		return nil, domainErrors.GetTokenValidationError(domainErrors.TokenMalformedErrorCode)
	}
	output := &domainModel.TokenUserJwtOutput{
		UserId:    out.UserId,
		TokenId:   out.ID,
		Role:      out.Role,
		Issuer:    out.Issuer,
		Subject:   out.Subject,
		Audience:  out.Audience,
		ExpiresAt: out.ExpiresAt.Time,
		IssuedAt:  out.IssuedAt.Time,
		NotBefore: out.NotBefore.Time,
	}
	return output, nil
}

/**
 * New token implementation
 */
func NewTokenService(
	secret string,
	issuer string,
	subject string,
	audience []string,
) domainToken.ITokenService {
	return &TokenService{
		secret:   secret,
		issuer:   issuer,
		subject:  subject,
		audience: audience,
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
