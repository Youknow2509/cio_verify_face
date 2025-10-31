// src/services/api/attendance.api.ts
/**
 * Attendance Service APIs
 * Quản lý chấm công: check-in, check-out, lịch sử
 */

import { http } from '../http';
import type { ApiResponse, AttendanceRecord, PaginatedResponse } from '@/types';

const API_PREFIX = '/api/v1/attendance';

/**
 * Check-in via face recognition
 * @param faceImage Base64 encoded face image
 * @returns Check-in result
 */
export async function checkInAPI(faceImage: string): Promise<ApiResponse<AttendanceRecord>> {
  return http.post(`${API_PREFIX}/check-in`, { faceImage });
}

/**
 * Check-out via face recognition
 * @param faceImage Base64 encoded face image
 * @returns Check-out result
 */
export async function checkOutAPI(faceImage: string): Promise<ApiResponse<AttendanceRecord>> {
  return http.post(`${API_PREFIX}/check-out`, { faceImage });
}

/**
 * Get attendance records with filters
 * @param startDate Start date (YYYY-MM-DD)
 * @param endDate End date (YYYY-MM-DD)
 * @param userId User/Employee ID (optional)
 * @param page Page number
 * @param limit Items per page
 * @returns Paginated list of attendance records
 */
export async function getAttendanceRecordsAPI(
  startDate?: string,
  endDate?: string,
  userId?: string,
  page?: number,
  limit?: number
): Promise<ApiResponse<PaginatedResponse<AttendanceRecord>>> {
  const params = new URLSearchParams();
  
  if (startDate) params.append('startDate', startDate);
  if (endDate) params.append('endDate', endDate);
  if (userId) params.append('userId', userId);
  if (page) params.append('page', page.toString());
  if (limit) params.append('limit', limit.toString());

  const queryString = params.toString();
  const url = queryString ? `${API_PREFIX}/records?${queryString}` : `${API_PREFIX}/records`;
  
  return http.get(url);
}

/**
 * Get single attendance record by ID
 * @param recordId Attendance record ID
 * @returns Attendance record data
 */
export async function getAttendanceRecordAPI(recordId: string): Promise<ApiResponse<AttendanceRecord>> {
  return http.get(`${API_PREFIX}/records/${recordId}`);
}

/**
 * Get current user's attendance history
 * @param startDate Start date (YYYY-MM-DD) optional
 * @param endDate End date (YYYY-MM-DD) optional
 * @param page Page number
 * @param limit Items per page
 * @returns Paginated list of personal attendance records
 */
export async function getMyAttendanceHistoryAPI(
  startDate?: string,
  endDate?: string,
  page?: number,
  limit?: number
): Promise<ApiResponse<PaginatedResponse<AttendanceRecord>>> {
  const params = new URLSearchParams();
  
  if (startDate) params.append('startDate', startDate);
  if (endDate) params.append('endDate', endDate);
  if (page) params.append('page', page.toString());
  if (limit) params.append('limit', limit.toString());

  const queryString = params.toString();
  const url = queryString 
    ? `${API_PREFIX}/history/my?${queryString}` 
    : `${API_PREFIX}/history/my`;
  
  return http.get(url);
}
