/**
 * Services Module
 * 
 * Main entry point for all services:
 * - API services (src/services/api/)
 * - HTTP client (src/services/http.ts)
 * - Error handling (src/services/error-handler.ts)
 * - Mock services (src/services/mock/) - for fallback/testing
 */

// === API Services ===
// These are the main services to use for API calls

export * from '@/services/api';
export { http, setAuthToken, clearAuthToken } from '@/services/http';
export type { RequestConfig, HttpError } from '@/services/http';

// === Error Handling ===
export {
  handleApiError,
  isAuthError,
  isValidationError,
  isNotFoundError,
  isServerError,
  getUserFriendlyMessage,
  type ErrorInfo
} from '@/services/error-handler';

// === HTTP Interceptor ===
export {
  configureHttpInterceptor,
  logRequest,
  logResponse,
  logError,
  isRetryableError,
  getRetryDelay,
  sleep,
  type InterceptorConfig
} from '@/services/http-interceptor';

// === Mock Services (Legacy) ===
// Import from these only for testing/fallback scenarios
// Note: Mock services are deprecated, use API services above
export * as mockServices from './mock';