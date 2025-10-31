// src/services/api/account.api.ts
/**
 * Account Management APIs
 * Quản lý tài khoản người dùng cá nhân
 * 
 * Note: User profile endpoints are part of Identity & Organization Service (/api/v1/users)
 * This file handles personal account operations like updating profile, changing password
 */

import { http } from '../http';
import type { ApiResponse, AccountProfile, ChangePasswordPayload } from '@/types';

// Get current user profile from Auth Service
const AUTH_PREFIX = '/api/v1/auth';
// Update user profile via Users Service
const USERS_PREFIX = '/api/v1/users';

/**
 * Get current user account profile
 * (This calls the auth /me endpoint)
 * @returns Current user profile data
 */
export async function getAccountProfileAPI(): Promise<ApiResponse<AccountProfile>> {
  return http.get(`${AUTH_PREFIX}/me`);
}

/**
 * Update own account profile
 * @param data Updated profile data (name, email, etc.)
 * @returns Updated account profile
 */
export async function updateAccountProfileAPI(
  data: Partial<Omit<AccountProfile, 'id'>>
): Promise<ApiResponse<AccountProfile>> {
  // Get current user first to get their ID
  const profileRes = await http.get<ApiResponse<AccountProfile>>(`${AUTH_PREFIX}/me`);
  const userId = (profileRes as any)?.data?.id;
  
  if (!userId) {
    throw new Error('Cannot update profile: User ID not found');
  }
  
  return http.put(`${USERS_PREFIX}/${userId}`, data);
}

/**
 * Change user password
 * @param payload Password change payload (oldPassword, newPassword)
 * @returns Void
 */
export async function changeAccountPasswordAPI(
  payload: ChangePasswordPayload
): Promise<ApiResponse<void>> {
  return http.post(`${AUTH_PREFIX}/change-password`, payload);
}

/**
 * Upload account avatar
 * @param file Avatar image file
 * @returns Updated account profile with new avatar URL
 */
export async function uploadAccountAvatarAPI(file: File): Promise<ApiResponse<AccountProfile>> {
  const formData = new FormData();
  formData.append('file', file);
  
  // Get current user ID first
  const profileRes = await http.get<ApiResponse<AccountProfile>>(`${AUTH_PREFIX}/me`);
  const userId = (profileRes as any)?.data?.id;
  
  if (!userId) {
    throw new Error('Cannot upload avatar: User ID not found');
  }
  
  return http.post(`${USERS_PREFIX}/${userId}/avatar`, formData);
}
