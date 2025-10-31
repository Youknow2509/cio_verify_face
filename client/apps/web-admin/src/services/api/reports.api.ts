// src/services/api/reports.api.ts
/**
 * Analytics & Reporting Service APIs
 * Báo cáo thống kê chấm công, quản lý
 */

import { http } from '../http';
import type { ApiResponse } from '@/types';

const API_PREFIX = '/api/v1/reports';

export interface DailyReport {
  date: string;
  totalCheckIns: number;
  totalCheckOuts: number;
  lateArrivals: number;
  absentEmployees: number;
  presentEmployees: number;
  averageCheckInTime?: string;
}

export interface SummaryReport {
  period: string;
  totalAttendance: number;
  averageAttendanceRate: number;
  totalLateArrivals: number;
  departmentBreakdown: {
    department: string;
    attendance: number;
    attendanceRate: number;
  }[];
  topAbsentees: {
    userId: string;
    userName: string;
    absenceCount: number;
  }[];
  topLateArrivals: {
    userId: string;
    userName: string;
    lateCount: number;
  }[];
}

export interface ExportReportParams {
  startDate?: string;
  endDate?: string;
  format?: 'csv' | 'pdf' | 'excel';
  department?: string;
}

/**
 * Get daily attendance report
 * @param date Report date (YYYY-MM-DD)
 * @param department Optional department filter
 * @returns Daily report
 */
export async function getDailyReportAPI(
  date?: string,
  department?: string
): Promise<ApiResponse<DailyReport>> {
  const params = new URLSearchParams();
  
  if (date) params.append('date', date);
  if (department) params.append('department', department);

  const queryString = params.toString();
  const url = queryString ? `${API_PREFIX}/daily?${queryString}` : `${API_PREFIX}/daily`;
  
  return http.get(url);
}

/**
 * Get summary report for a period
 * @param startDate Start date (YYYY-MM-DD)
 * @param endDate End date (YYYY-MM-DD)
 * @param department Optional department filter
 * @returns Summary report
 */
export async function getSummaryReportAPI(
  startDate?: string,
  endDate?: string,
  department?: string
): Promise<ApiResponse<SummaryReport>> {
  const params = new URLSearchParams();
  
  if (startDate) params.append('startDate', startDate);
  if (endDate) params.append('endDate', endDate);
  if (department) params.append('department', department);

  const queryString = params.toString();
  const url = queryString ? `${API_PREFIX}/summary?${queryString}` : `${API_PREFIX}/summary`;
  
  return http.get(url);
}

/**
 * Export report in various formats
 * @param params Export parameters (format, date range, etc.)
 * @returns Report file URL or binary data
 */
export async function exportReportAPI(
  params: ExportReportParams
): Promise<ApiResponse<{ url: string }>> {
  const searchParams = new URLSearchParams();
  
  if (params.startDate) searchParams.append('startDate', params.startDate);
  if (params.endDate) searchParams.append('endDate', params.endDate);
  if (params.format) searchParams.append('format', params.format);
  if (params.department) searchParams.append('department', params.department);

  const queryString = searchParams.toString();
  const url = queryString ? `${API_PREFIX}/export?${queryString}` : `${API_PREFIX}/export`;
  
  return http.get(url);
}
