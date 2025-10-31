// src/config/api.config.ts
/**
 * API Configuration
 * 
 * Centralized configuration for API client
 */

export interface ApiConfig {
  baseUrl: string;
  timeout: number;
  retryAttempts: number;
  retryDelay: number;
}

/**
 * Get API configuration based on environment
 */
export function getApiConfig(): ApiConfig {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
  const timeout = parseInt(import.meta.env.VITE_API_TIMEOUT || '30000', 10);
  const retryAttempts = parseInt(import.meta.env.VITE_API_RETRY_ATTEMPTS || '3', 10);
  const retryDelay = parseInt(import.meta.env.VITE_API_RETRY_DELAY || '1000', 10);

  return {
    baseUrl,
    timeout,
    retryAttempts,
    retryDelay
  };
}

/**
 * Development API configuration
 */
export const DEV_API_CONFIG: ApiConfig = {
  baseUrl: 'http://localhost:8080/api',
  timeout: 30000,
  retryAttempts: 3,
  retryDelay: 1000
};

/**
 * Production API configuration
 */
export const PROD_API_CONFIG: ApiConfig = {
  baseUrl: '/api',
  timeout: 30000,
  retryAttempts: 2,
  retryDelay: 2000
};

/**
 * Get environment-specific API configuration
 */
export function getEnvironmentApiConfig(): ApiConfig {
  const isDev = import.meta.env.DEV;
  return isDev ? DEV_API_CONFIG : PROD_API_CONFIG;
}
