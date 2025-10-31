// src/services/api/shifts.api.ts
/**
 * Workforce Service APIs
 * Quản lý ca làm việc (shifts) và lịch làm việc (schedules)
 */

import { http } from '../http';
import type { ApiResponse, Shift, Schedule, PaginatedResponse } from '@/types';

// ===== Shifts Management =====
const SHIFTS_PREFIX = '/api/v1/shifts';

/**
 * Get list of shifts
 * @param page Page number
 * @param limit Items per page
 * @param active Filter by active status
 * @returns Paginated list of shifts
 */
export async function getShiftsAPI(
  page?: number,
  limit?: number,
  active?: boolean
): Promise<ApiResponse<PaginatedResponse<Shift>>> {
  const params = new URLSearchParams();
  
  if (page) params.append('page', page.toString());
  if (limit) params.append('limit', limit.toString());
  if (active !== undefined) params.append('active', active.toString());

  const queryString = params.toString();
  const url = queryString ? `${SHIFTS_PREFIX}?${queryString}` : SHIFTS_PREFIX;
  
  return http.get(url);
}

/**
 * Get shift by ID
 * @param shiftId Shift ID
 * @returns Shift data
 */
export async function getShiftAPI(shiftId: string): Promise<ApiResponse<Shift>> {
  return http.get(`${SHIFTS_PREFIX}/${shiftId}`);
}

/**
 * Create new shift
 * @param data Shift data
 * @returns Created shift
 */
export async function createShiftAPI(
  data: Omit<Shift, 'id' | 'createdAt'>
): Promise<ApiResponse<Shift>> {
  return http.post(SHIFTS_PREFIX, data);
}

/**
 * Update shift
 * @param shiftId Shift ID
 * @param data Updated shift data
 * @returns Updated shift
 */
export async function updateShiftAPI(
  shiftId: string,
  data: Partial<Omit<Shift, 'id' | 'createdAt'>>
): Promise<ApiResponse<Shift>> {
  return http.put(`${SHIFTS_PREFIX}/${shiftId}`, data);
}

/**
 * Delete shift
 * @param shiftId Shift ID
 * @returns Void
 */
export async function deleteShiftAPI(shiftId: string): Promise<ApiResponse<void>> {
  return http.delete(`${SHIFTS_PREFIX}/${shiftId}`);
}

// ===== Schedules Management =====
const SCHEDULES_PREFIX = '/api/v1/schedules';

/**
 * Get list of schedules
 * @param page Page number
 * @param limit Items per page
 * @param userId Filter by user/employee ID
 * @returns Paginated list of schedules
 */
export async function getSchedulesAPI(
  page?: number,
  limit?: number,
  userId?: string
): Promise<ApiResponse<PaginatedResponse<Schedule>>> {
  const params = new URLSearchParams();
  
  if (page) params.append('page', page.toString());
  if (limit) params.append('limit', limit.toString());
  if (userId) params.append('userId', userId);

  const queryString = params.toString();
  const url = queryString ? `${SCHEDULES_PREFIX}?${queryString}` : SCHEDULES_PREFIX;
  
  return http.get(url);
}

/**
 * Get schedule by ID
 * @param scheduleId Schedule ID
 * @returns Schedule data
 */
export async function getScheduleAPI(scheduleId: string): Promise<ApiResponse<Schedule>> {
  return http.get(`${SCHEDULES_PREFIX}/${scheduleId}`);
}

/**
 * Create new schedule
 * @param data Schedule data (userId, shiftId, startDate, endDate)
 * @returns Created schedule
 */
export async function createScheduleAPI(
  data: Omit<Schedule, 'id' | 'createdAt'>
): Promise<ApiResponse<Schedule>> {
  return http.post(SCHEDULES_PREFIX, data);
}

/**
 * Update schedule
 * @param scheduleId Schedule ID
 * @param data Updated schedule data
 * @returns Updated schedule
 */
export async function updateScheduleAPI(
  scheduleId: string,
  data: Partial<Omit<Schedule, 'id' | 'createdAt'>>
): Promise<ApiResponse<Schedule>> {
  return http.put(`${SCHEDULES_PREFIX}/${scheduleId}`, data);
}

/**
 * Delete schedule
 * @param scheduleId Schedule ID
 * @returns Void
 */
export async function deleteScheduleAPI(scheduleId: string): Promise<ApiResponse<void>> {
  return http.delete(`${SCHEDULES_PREFIX}/${scheduleId}`);
}
