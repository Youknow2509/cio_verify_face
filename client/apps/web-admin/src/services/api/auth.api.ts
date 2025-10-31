// src/services/api/auth.api.ts
/**
 * Auth Service APIs
 * Quản lý xác thực người dùng & thiết bị
 */

import { http } from '../http';
import type { ApiResponse, User } from '@/types';

const API_PREFIX = '/api/v1/auth';

/**
 * Đăng nhập tài khoản (user/device)
 * @param email Email hoặc username
 * @param password Password
 * @returns Access token và thông tin user
 */
export async function loginAPI(
  email: string,
  password: string
): Promise<ApiResponse<{ token: string; user: User }>> {
  return http.post(`${API_PREFIX}/login`, { email, password });
}

/**
 * Đăng xuất tài khoản
 * @returns Void
 */
export async function logoutAPI(): Promise<ApiResponse<void>> {
  return http.post(`${API_PREFIX}/logout`, {});
}

/**
 * Làm mới access token
 * @returns New access token
 */
export async function refreshTokenAPI(): Promise<ApiResponse<{ token: string }>> {
  return http.post(`${API_PREFIX}/refresh`, {});
}

/**
 * Kích hoạt thiết bị chấm công
 * @param deviceCode Mã thiết bị
 * @param deviceSecret Secret key thiết bị
 * @returns Kết quả kích hoạt
 */
export async function activateDeviceAPI(
  deviceCode: string,
  deviceSecret: string
): Promise<ApiResponse<{ deviceId: string; activated: boolean }>> {
  return http.post(`${API_PREFIX}/device/activate`, {
    deviceCode,
    deviceSecret
  });
}

/**
 * Lấy thông tin tài khoản hiện tại
 * @returns Thông tin user đang đăng nhập
 */
export async function getCurrentUserAPI(): Promise<ApiResponse<User>> {
  return http.get(`${API_PREFIX}/me`);
}
