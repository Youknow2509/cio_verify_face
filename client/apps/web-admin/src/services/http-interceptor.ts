// src/services/http-interceptor.ts
/**
 * HTTP Request Interceptor
 * 
 * Handles request/response interception for features like:
 * - Auth token management
 * - Request/response logging
 * - Error handling
 * - Retry logic
 */

import { handleApiError } from './error-handler';

export interface InterceptorConfig {
  enableLogging?: boolean;
  enableRetry?: boolean;
  maxRetries?: number;
  retryDelay?: number;
}

const DEFAULT_CONFIG: InterceptorConfig = {
  enableLogging: true,
  enableRetry: true,
  maxRetries: 3,
  retryDelay: 1000
};

let interceptorConfig = DEFAULT_CONFIG;

/**
 * Configure HTTP interceptor
 * @param config Interceptor configuration
 */
export function configureHttpInterceptor(config: Partial<InterceptorConfig>) {
  interceptorConfig = { ...interceptorConfig, ...config };
}

/**
 * Log HTTP request
 * @param method HTTP method
 * @param url Request URL
 * @param body Request body
 */
export function logRequest(method: string, url: string, body?: any) {
  if (!interceptorConfig.enableLogging) return;

  console.group(`📤 ${method} ${url}`);
  if (body) {
    console.log('Body:', body);
  }
  console.groupEnd();
}

/**
 * Log HTTP response
 * @param method HTTP method
 * @param url Request URL
 * @param status Response status
 * @param data Response data
 */
export function logResponse(method: string, url: string, status: number, data?: any) {
  if (!interceptorConfig.enableLogging) return;

  const statusEmoji = status >= 200 && status < 300 ? '✅' : '❌';
  console.group(`${statusEmoji} ${method} ${url} [${status}]`);
  if (data) {
    console.log('Response:', data);
  }
  console.groupEnd();
}

/**
 * Log HTTP error
 * @param method HTTP method
 * @param url Request URL
 * @param error Error object
 */
export function logError(method: string, url: string, error: unknown) {
  if (!interceptorConfig.enableLogging) return;

  const errorInfo = handleApiError(error);
  console.group(`❌ ${method} ${url} [${errorInfo.code}]`);
  console.error('Error:', errorInfo);
  console.groupEnd();
}

/**
 * Check if error is retryable
 * @param error Error object
 * @returns True if error can be retried
 */
export function isRetryableError(error: unknown): boolean {
  if (!interceptorConfig.enableRetry) return false;

  // Retry on network errors
  if (error instanceof Error) {
    if (error.message.includes('Network') || error.message.includes('timeout')) {
      return true;
    }
    
    // Retry on specific HTTP status codes
    if ('status' in error && typeof error.status === 'number') {
      const status = error.status;
      // Retry on 408 (Timeout), 429 (Too Many Requests), 5xx (Server errors)
      return status === 408 || status === 429 || (status >= 500 && status < 600);
    }
  }

  return false;
}

/**
 * Get retry delay for attempt number
 * @param attemptNumber Current attempt number (0-based)
 * @returns Delay in milliseconds
 */
export function getRetryDelay(attemptNumber: number): number {
  const baseDelay = interceptorConfig.retryDelay || 1000;
  // Exponential backoff: 1s, 2s, 4s, 8s, etc.
  return baseDelay * Math.pow(2, attemptNumber);
}

/**
 * Sleep for specified duration
 * @param ms Duration in milliseconds
 */
export function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}
