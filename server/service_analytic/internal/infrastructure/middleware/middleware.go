package middleware

// Middleware package contains authentication and session management middleware
// for both HTTP (Gin) and gRPC servers.
//
// HTTP Middleware:
// - HTTPAuthMiddleware: Validates JWT tokens via service_auth and sets session in context
//
// gRPC Middleware:
// - SessionInterceptor: Extracts pre-authenticated session from metadata (inter-service)
