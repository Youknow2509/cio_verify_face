// src/services/api/employees.api.ts
/**
 * Identity & Organization Service APIs
 * Quản lý công ty, user, nhân viên, dữ liệu khuôn mặt
 */

import { http } from '../http';
import type { ApiResponse, Employee, FaceData, PaginatedResponse, EmployeeFilter } from '@/types';

const API_PREFIX = '/api/v1/users';

/**
 * Lấy danh sách nhân viên/user trong công ty
 * @param filter Filter options (page, limit, search, department, etc.)
 * @returns Paginated list of users
 */
export async function getEmployeesAPI(
  filter?: EmployeeFilter
): Promise<ApiResponse<PaginatedResponse<Employee>>> {
  const params = new URLSearchParams();
  
  if (filter) {
    if (filter.page) params.append('page', filter.page.toString());
    if (filter.limit) params.append('limit', filter.limit.toString());
    if (filter.search) params.append('search', filter.search);
    if (filter.department) params.append('department', filter.department);
    if (filter.active !== undefined) params.append('active', filter.active.toString());
    if (filter.sortBy) params.append('sortBy', filter.sortBy);
    if (filter.sortOrder) params.append('sortOrder', filter.sortOrder);
  }

  const queryString = params.toString();
  const url = queryString ? `${API_PREFIX}?${queryString}` : API_PREFIX;
  
  return http.get<ApiResponse<PaginatedResponse<Employee>>>(url);
}

/**
 * Xem thông tin chi tiết nhân viên
 * @param userId User ID
 * @returns Employee data
 */
export async function getEmployeeAPI(userId: string): Promise<ApiResponse<Employee>> {
  return http.get(`${API_PREFIX}/${userId}`);
}

/**
 * Thêm mới nhân viên
 * @param data Employee data
 * @returns Created employee
 */
export async function createEmployeeAPI(
  data: Omit<Employee, 'id' | 'createdAt' | 'updatedAt'>
): Promise<ApiResponse<Employee>> {
  return http.post(API_PREFIX, data);
}

/**
 * Sửa thông tin nhân viên
 * @param userId User ID
 * @param data Updated employee data
 * @returns Updated employee
 */
export async function updateEmployeeAPI(
  userId: string,
  data: Partial<Omit<Employee, 'id' | 'createdAt' | 'updatedAt'>>
): Promise<ApiResponse<Employee>> {
  return http.put(`${API_PREFIX}/${userId}`, data);
}

/**
 * Vô hiệu hóa/xóa nhân viên
 * @param userId User ID
 * @returns Void
 */
export async function deleteEmployeeAPI(userId: string): Promise<ApiResponse<void>> {
  return http.delete(`${API_PREFIX}/${userId}`);
}

/**
 * Lấy danh sách ảnh khuôn mặt của nhân viên
 * @param userId User ID
 * @returns List of face data
 */
export async function getEmployeeFaceDataAPI(userId: string): Promise<ApiResponse<FaceData[]>> {
  return http.get(`${API_PREFIX}/${userId}/face-data`);
}

/**
 * Đăng ký/upload ảnh khuôn mặt mới
 * @param userId User ID
 * @param file Image file
 * @returns Uploaded face data
 */
export async function uploadFaceDataAPI(
  userId: string,
  file: File
): Promise<ApiResponse<FaceData>> {
  const formData = new FormData();
  formData.append('file', file);
  
  return http.post(`${API_PREFIX}/${userId}/face-data`, formData);
}

/**
 * Xóa ảnh khuôn mặt
 * @param userId User ID
 * @param faceId Face ID
 * @returns Void
 */
export async function deleteFaceDataAPI(
  userId: string,
  faceId: string
): Promise<ApiResponse<void>> {
  return http.delete(`${API_PREFIX}/${userId}/face-data/${faceId}`);
}
