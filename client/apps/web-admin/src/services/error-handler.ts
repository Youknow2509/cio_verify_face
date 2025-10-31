// src/services/error-handler.ts
/**
 * API Error Handler
 * 
 * Centralized error handling for API requests
 */

import type { HttpError } from './http';

export interface ErrorInfo {
  code: string;
  message: string;
  statusCode?: number;
  details?: any;
}

/**
 * Handle API errors
 * @param error Error from API request
 * @returns Formatted error info
 */
export function handleApiError(error: unknown): ErrorInfo {
  if (error instanceof Error) {
    // Check if it's HttpError
    if ('status' in error && typeof error.status === 'number') {
      const httpError = error as HttpError;
      return {
        code: `HTTP_${httpError.status}`,
        message: httpError.message || httpError.statusText,
        statusCode: httpError.status,
        details: null
      };
    }

    // Network or timeout errors
    if (error.message.includes('timeout')) {
      return {
        code: 'TIMEOUT',
        message: 'Request timeout. Please try again.',
        statusCode: 0
      };
    }

    if (error.message.includes('Network')) {
      return {
        code: 'NETWORK_ERROR',
        message: 'Network error. Please check your connection.',
        statusCode: 0
      };
    }

    return {
      code: 'UNKNOWN_ERROR',
      message: error.message || 'An unknown error occurred',
      statusCode: 0
    };
  }

  return {
    code: 'UNKNOWN_ERROR',
    message: 'An unknown error occurred',
    statusCode: 0
  };
}

/**
 * Check if error is authentication error
 * @param error Error info
 * @returns True if error is auth-related
 */
export function isAuthError(error: ErrorInfo): boolean {
  return error.statusCode === 401 || error.statusCode === 403;
}

/**
 * Check if error is validation error
 * @param error Error info
 * @returns True if error is validation-related
 */
export function isValidationError(error: ErrorInfo): boolean {
  return error.statusCode === 400 || error.statusCode === 422;
}

/**
 * Check if error is not found error
 * @param error Error info
 * @returns True if error is not found
 */
export function isNotFoundError(error: ErrorInfo): boolean {
  return error.statusCode === 404;
}

/**
 * Check if error is server error
 * @param error Error info
 * @returns True if error is server-related
 */
export function isServerError(error: ErrorInfo): boolean {
  return error.statusCode ? error.statusCode >= 500 : false;
}

/**
 * Get user-friendly error message
 * @param error Error info
 * @returns User-friendly message
 */
export function getUserFriendlyMessage(error: ErrorInfo): string {
  if (isAuthError(error)) {
    return 'Authentication failed. Please login again.';
  }

  if (isValidationError(error)) {
    return 'Invalid data. Please check your input.';
  }

  if (isNotFoundError(error)) {
    return 'Resource not found.';
  }

  if (isServerError(error)) {
    return 'Server error. Please try again later.';
  }

  if (error.code === 'TIMEOUT') {
    return 'Request timeout. Please try again.';
  }

  if (error.code === 'NETWORK_ERROR') {
    return 'Network connection error. Please check your internet.';
  }

  return error.message || 'An error occurred. Please try again.';
}
