// src/services/api-helpers.ts
/**
 * API Helper Utilities
 * 
 * Utilities for working with API services
 */

import type { Employee, Device, Shift, AttendanceRecord } from '@/types';

/**
 * Build query string from object
 * @param params Query parameters object
 * @returns Query string (without leading ?)
 */
export function buildQueryString(params: Record<string, any>): string {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value !== null && value !== undefined && value !== '') {
      searchParams.append(key, String(value));
    }
  });

  return searchParams.toString();
}

/**
 * Create employee payload from form data
 * @param formData Form data
 * @returns Employee creation payload
 */
export function createEmployeePayload(formData: any): Omit<Employee, 'id' | 'createdAt' | 'updatedAt'> {
  return {
    code: formData.code,
    name: formData.name,
    email: formData.email,
    department: formData.department,
    position: formData.position,
    active: formData.active ?? true,
    faceCount: 0
  };
}

/**
 * Create device payload from form data
 * @param formData Form data
 * @returns Device creation payload
 */
export function createDevicePayload(formData: any): Omit<Device, 'id' | 'createdAt'> {
  return {
    name: formData.name,
    location: formData.location,
    status: formData.status ?? 'offline',
    model: formData.model,
    ipAddress: formData.ipAddress
  };
}

/**
 * Create shift payload from form data
 * @param formData Form data
 * @returns Shift creation payload
 */
export function createShiftPayload(formData: any): Omit<Shift, 'id' | 'createdAt'> {
  return {
    name: formData.name,
    start: formData.start,
    end: formData.end,
    description: formData.description,
    active: formData.active ?? true
  };
}

/**
 * Create attendance record payload from form data
 * @param formData Form data
 * @returns Attendance record creation payload
 */
export function createAttendancePayload(formData: any): Omit<AttendanceRecord, 'id'> {
  return {
    employeeId: formData.employeeId,
    employeeName: formData.employeeName,
    date: formData.date,
    checkIn: formData.checkIn,
    checkOut: formData.checkOut,
    isLate: formData.isLate ?? false,
    shiftId: formData.shiftId,
    totalHours: formData.totalHours,
    deviceId: formData.deviceId
  };
}

/**
 * Format date for API (YYYY-MM-DD)
 * @param date Date object or string
 * @returns Formatted date string
 */
export function formatDateForAPI(date: Date | string): string {
  if (typeof date === 'string') return date;
  
  return date.toISOString().split('T')[0];
}

/**
 * Format datetime for API (ISO 8601)
 * @param date Date object or string
 * @returns ISO formatted datetime
 */
export function formatDatetimeForAPI(date: Date | string): string {
  if (typeof date === 'string') return date;
  
  return date.toISOString();
}

/**
 * Parse API date response (YYYY-MM-DD or ISO 8601)
 * @param dateString Date string from API
 * @returns Date object
 */
export function parseAPIDate(dateString: string): Date {
  return new Date(dateString);
}

/**
 * Check if form data is valid before API call
 * @param data Form data
 * @param requiredFields Required field names
 * @returns Array of missing fields or empty if valid
 */
export function validateFormData(data: any, requiredFields: string[]): string[] {
  const missing: string[] = [];

  requiredFields.forEach(field => {
    const value = data[field];
    if (value === null || value === undefined || value === '') {
      missing.push(field);
    }
  });

  return missing;
}

/**
 * Merge pagination defaults with provided options
 * @param page Page number (1-based)
 * @param limit Items per page
 * @returns Pagination object
 */
export function getPaginationParams(page?: number, limit?: number) {
  return {
    page: page || 1,
    limit: limit || 10
  };
}

/**
 * Extract pagination info from API response
 * @param data Response data with pagination
 * @returns Pagination info
 */
export function extractPaginationInfo(data: any) {
  return {
    page: data.page || 1,
    limit: data.limit || 10,
    total: data.total || 0,
    totalPages: Math.ceil((data.total || 0) / (data.limit || 10))
  };
}

/**
 * Create form data object for file uploads
 * @param file File to upload
 * @param additionalData Additional form data
 * @returns FormData object
 */
export function createFormData(file: File, additionalData?: Record<string, any>): FormData {
  const formData = new FormData();
  formData.append('file', file);

  if (additionalData) {
    Object.entries(additionalData).forEach(([key, value]) => {
      if (value !== null && value !== undefined) {
        formData.append(key, value instanceof File ? value : String(value));
      }
    });
  }

  return formData;
}

/**
 * Build filter query from filter object
 * @param filter Filter object
 * @returns URL query string
 */
export function buildFilterQuery(filter: Record<string, any>): URLSearchParams {
  const params = new URLSearchParams();

  Object.entries(filter).forEach(([key, value]) => {
    if (value !== null && value !== undefined && value !== '') {
      params.append(key, String(value));
    }
  });

  return params;
}

/**
 * Get file extension from filename
 * @param filename File name
 * @returns File extension
 */
export function getFileExtension(filename: string): string {
  const parts = filename.split('.');
  return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : '';
}

/**
 * Check if file is valid image
 * @param file File object
 * @returns True if valid image
 */
export function isValidImageFile(file: File): boolean {
  const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
  const validExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp'];
  
  const isValidType = validTypes.includes(file.type);
  const isValidExtension = validExtensions.includes(getFileExtension(file.name));
  
  return isValidType && isValidExtension;
}

/**
 * Check if file size is valid (in MB)
 * @param file File object
 * @param maxSizeMB Max file size in MB
 * @returns True if valid
 */
export function isValidFileSize(file: File, maxSizeMB: number = 10): boolean {
  return file.size <= maxSizeMB * 1024 * 1024;
}
