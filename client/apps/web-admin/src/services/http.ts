// src/services/http.ts

export interface RequestConfig {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  headers?: Record<string, string>;
  body?: any;
  timeout?: number;
}

export class HttpError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    message?: string
  ) {
    super(message || `HTTP ${status}: ${statusText}`);
    this.name = 'HttpError';
  }
}

const DEFAULT_TIMEOUT = 10000;
const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

function getAuthToken(): string | null {
  try {
    return localStorage.getItem('auth_token');
  } catch {
    return null;
  }
}

function createRequestConfig(config: RequestConfig = {}): RequestInit {
  const {
    method = 'GET',
    headers = {},
    body,
    timeout = DEFAULT_TIMEOUT
  } = config;

  const requestHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    ...headers
  };

  // Add auth token if available
  const token = getAuthToken();
  if (token) {
    requestHeaders.Authorization = `Bearer ${token}`;
  }

  const requestConfig: RequestInit = {
    method,
    headers: requestHeaders,
    signal: AbortSignal.timeout(timeout)
  };

  if (body && method !== 'GET') {
    requestConfig.body = typeof body === 'string' ? body : JSON.stringify(body);
  }

  return requestConfig;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = response.statusText;
    
    try {
      const errorData = await response.json();
      errorMessage = errorData.message || errorData.error || errorMessage;
    } catch {
      // If response is not JSON, use statusText
    }
    
    throw new HttpError(response.status, response.statusText, errorMessage);
  }

  try {
    return await response.json();
  } catch {
    // If response is not JSON (e.g., 204 No Content), return empty object
    return {} as T;
  }
}

export async function httpRequest<T>(
  endpoint: string,
  config: RequestConfig = {}
): Promise<T> {
  const url = endpoint.startsWith('http') ? endpoint : `${BASE_URL}${endpoint}`;
  const requestConfig = createRequestConfig(config);

  try {
    const response = await fetch(url, requestConfig);
    return await handleResponse<T>(response);
  } catch (error) {
    if (error instanceof HttpError) {
      throw error;
    }
    
    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw new Error('Request timeout');
      }
      throw new Error(`Network error: ${error.message}`);
    }
    
    throw new Error('Unknown error occurred');
  }
}

// Convenience methods
export const http = {
  get: <T>(endpoint: string, config?: Omit<RequestConfig, 'method' | 'body'>) =>
    httpRequest<T>(endpoint, { ...config, method: 'GET' }),

  post: <T>(endpoint: string, body?: any, config?: Omit<RequestConfig, 'method'>) =>
    httpRequest<T>(endpoint, { ...config, method: 'POST', body }),

  put: <T>(endpoint: string, body?: any, config?: Omit<RequestConfig, 'method'>) =>
    httpRequest<T>(endpoint, { ...config, method: 'PUT', body }),

  patch: <T>(endpoint: string, body?: any, config?: Omit<RequestConfig, 'method'>) =>
    httpRequest<T>(endpoint, { ...config, method: 'PATCH', body }),

  delete: <T>(endpoint: string, config?: Omit<RequestConfig, 'method' | 'body'>) =>
    httpRequest<T>(endpoint, { ...config, method: 'DELETE' })
};