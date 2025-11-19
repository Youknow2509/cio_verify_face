package middleware

import (
	"context"
	"fmt"
	"strconv"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// SessionInterceptor creates a gRPC unary interceptor for session extraction
// For inter-service gRPC calls, session info is passed directly in metadata (already authenticated)
func SessionInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		// Extract session info from metadata (set by calling service after auth)
		sessionInfo, err := extractSessionFromMetadata(md)
		if err != nil {
			global.Logger.Error("Failed to extract session info", "error", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid session info: %v", err)
		}

		// Add session info to context
		ctx = context.WithValue(ctx, constants.ContextKeySessionInfo, sessionInfo)

		// Call the actual handler
		return handler(ctx, req)
	}
}

// extractSessionFromMetadata extracts session info from gRPC metadata
// This expects the calling service to have already authenticated and populated these fields
func extractSessionFromMetadata(md metadata.MD) (*applicationModel.SessionInfo, error) {
	session := &applicationModel.SessionInfo{}

	// Extract user_id
	if userIDs := md.Get("user_id"); len(userIDs) > 0 {
		session.UserID = userIDs[0]
	} else {
		return nil, fmt.Errorf("user_id not found in metadata")
	}

	// Extract role
	if roles := md.Get("role"); len(roles) > 0 {
		roleInt, err := strconv.ParseInt(roles[0], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid role format")
		}
		session.Role = int32(roleInt)
	} else {
		return nil, fmt.Errorf("role not found in metadata")
	}

	// Extract session_id (optional)
	if sessionIDs := md.Get("session_id"); len(sessionIDs) > 0 {
		session.SessionID = sessionIDs[0]
	}

	// Extract company_id (optional but recommended for managers)
	if companyIDs := md.Get("company_id"); len(companyIDs) > 0 {
		session.CompanyID = companyIDs[0]
	}

	// Extract client_ip (optional)
	if clientIPs := md.Get("client_ip"); len(clientIPs) > 0 {
		session.ClientIP = clientIPs[0]
	}

	// Extract client_agent (optional)
	if clientAgents := md.Get("client_agent"); len(clientAgents) > 0 {
		session.ClientAgent = clientAgents[0]
	}

	return session, nil
}

// GetSessionInfoFromContext extracts session info from gRPC context
func GetSessionInfoFromContext(ctx context.Context) (*applicationModel.SessionInfo, error) {
	sessionInfo, ok := ctx.Value(constants.ContextKeySessionInfo).(*applicationModel.SessionInfo)
	if !ok {
		return nil, fmt.Errorf("session info not found in context")
	}
	return sessionInfo, nil
}
