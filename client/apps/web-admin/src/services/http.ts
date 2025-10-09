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
const TOKEN_STORAGE_KEY = 'auth_token';

let inMemoryToken: string | undefined;

function readTokenFromStorage(): string | undefined {
  if (typeof window === 'undefined') {
    return undefined;
  }

  try {
    const localToken = window.localStorage?.getItem(TOKEN_STORAGE_KEY) ?? undefined;
    if (localToken) {
      return localToken;
    }
  } catch {
    // Ignore storage errors
  }

  try {
    const sessionToken = window.sessionStorage?.getItem(TOKEN_STORAGE_KEY) ?? undefined;
    if (sessionToken) {
      return sessionToken;
    }
  } catch {
    // Ignore storage errors
  }

  return undefined;
}

function persistToken(token?: string) {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    if (typeof token === 'string' && token.length > 0) {
      window.localStorage?.setItem(TOKEN_STORAGE_KEY, token);
    } else {
      window.localStorage?.removeItem(TOKEN_STORAGE_KEY);
    }
  } catch {
    // Ignore storage errors
  }

  try {
    if (typeof token === 'string' && token.length > 0) {
      window.sessionStorage?.setItem(TOKEN_STORAGE_KEY, token);
    } else {
      window.sessionStorage?.removeItem(TOKEN_STORAGE_KEY);
    }
  } catch {
    // Ignore storage errors
  }
}

export function setAuthToken(token?: string) {
  inMemoryToken = token && token.length > 0 ? token : undefined;
  persistToken(inMemoryToken);
}

export function clearAuthToken() {
  inMemoryToken = undefined;
  persistToken(undefined);
}

function getAuthToken(): string | null {
  if (inMemoryToken) {
    return inMemoryToken;
  }

  const storedToken = readTokenFromStorage();
  if (storedToken) {
    inMemoryToken = storedToken;
    return storedToken;
  }

  return null;
}

function createRequestConfig(config: RequestConfig = {}): RequestInit {
  const {
    method = 'GET',
    headers = {},
    body,
    timeout = DEFAULT_TIMEOUT
  } = config;

  const requestHeaders: Record<string, string> = {
    Accept: 'application/json',
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

  if (body !== undefined && method !== 'GET') {
    const isFormLike =
      body instanceof FormData ||
      body instanceof URLSearchParams ||
      body instanceof Blob ||
      body instanceof ArrayBuffer;

    if (isFormLike) {
      requestConfig.body = body as BodyInit;
    } else if (typeof body === 'string') {
      requestConfig.body = body;
      if (!requestHeaders['Content-Type']) {
        requestHeaders['Content-Type'] = 'text/plain;charset=UTF-8';
      }
    } else {
      requestConfig.body = JSON.stringify(body);
      if (!requestHeaders['Content-Type']) {
        requestHeaders['Content-Type'] = 'application/json';
      }
    }
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