// src/services/mock/auth.ts

import type { ApiResponse, User } from '../../types';

const MOCK_DELAY = 500;

const mockUser: User = {
  id: '1',
  email: 'admin@company.com',
  name: 'Nguyễn Văn Admin',
  role: 'CompanyAdmin',
  companyId: 'company-1',
  active: true
};

export async function login(email: string, password: string): Promise<ApiResponse<{ token: string; user: User }>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  // Mock validation
  if (email === 'admin@company.com' && password === 'password') {
    const token = 'mock-jwt-token-' + Date.now();
    
    // Store token in localStorage
    try {
      localStorage.setItem('auth_token', token);
      localStorage.setItem('user', JSON.stringify(mockUser));
    } catch {
      // Ignore storage errors
    }

    return {
      data: {
        token,
        user: mockUser
      }
    };
  }

  return {
    data: {
      token: '',
      user: {} as User
    },
    error: 'Email hoặc mật khẩu không đúng'
  };
}

export async function logout(): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  try {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
  } catch {
    // Ignore storage errors
  }

  return {
    data: undefined
  };
}

export async function refreshToken(): Promise<ApiResponse<{ token: string }>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  const newToken = 'mock-jwt-token-' + Date.now();
  
  try {
    localStorage.setItem('auth_token', newToken);
  } catch {
    // Ignore storage errors
  }

  return {
    data: {
      token: newToken
    }
  };
}

export async function getCurrentUser(): Promise<ApiResponse<User>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  try {
    const stored = localStorage.getItem('user');
    const user = stored ? JSON.parse(stored) : mockUser;
    
    return {
      data: user
    };
  } catch {
    return {
      data: mockUser
    };
  }
}

export async function changePassword(currentPassword: string, _newPassword: string): Promise<ApiResponse<void>> {
  await new Promise(resolve => setTimeout(resolve, MOCK_DELAY));

  // Mock validation
  if (currentPassword === 'password') {
    return {
      data: undefined
    };
  }

  return {
    data: undefined,
    error: 'Mật khẩu hiện tại không đúng'
  };
}