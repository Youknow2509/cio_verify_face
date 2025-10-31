// src/services/api/devices.api.ts
/**
 * Device Management Service APIs
 * Quản lý thiết bị nhân dạng khuôn mặt
 */

import { http } from '../http';
import type { ApiResponse, Device, PaginatedResponse } from '@/types';

const API_PREFIX = '/api/v1/devices';

/**
 * Get list of devices
 * @param page Page number
 * @param limit Items per page
 * @param search Search query
 * @returns Paginated list of devices
 */
export async function getDevicesAPI(
  page?: number,
  limit?: number,
  search?: string
): Promise<ApiResponse<PaginatedResponse<Device>>> {
  const params = new URLSearchParams();
  
  if (page) params.append('page', page.toString());
  if (limit) params.append('limit', limit.toString());
  if (search) params.append('search', search);

  const queryString = params.toString();
  const url = queryString ? `${API_PREFIX}?${queryString}` : API_PREFIX;
  
  return http.get(url);
}

/**
 * Get device by ID
 * @param id Device ID
 * @returns Device data
 */
export async function getDeviceAPI(id: string): Promise<ApiResponse<Device>> {
  return http.get(`${API_PREFIX}/${id}`);
}

/**
 * Create new device
 * @param data Device data
 * @returns Created device
 */
export async function createDeviceAPI(
  data: Omit<Device, 'id' | 'createdAt'>
): Promise<ApiResponse<Device>> {
  return http.post(API_PREFIX, data);
}

/**
 * Update device
 * @param id Device ID
 * @param data Updated device data
 * @returns Updated device
 */
export async function updateDeviceAPI(
  id: string,
  data: Partial<Omit<Device, 'id' | 'createdAt'>>
): Promise<ApiResponse<Device>> {
  return http.put(`${API_PREFIX}/${id}`, data);
}

/**
 * Delete device
 * @param id Device ID
 * @returns Void
 */
export async function deleteDeviceAPI(id: string): Promise<ApiResponse<void>> {
  return http.delete(`${API_PREFIX}/${id}`);
}

/**
 * Sync device (push config to device)
 * @param id Device ID
 * @returns Void
 */
export async function syncDeviceAPI(id: string): Promise<ApiResponse<void>> {
  return http.post(`${API_PREFIX}/${id}/sync`, {});
}
