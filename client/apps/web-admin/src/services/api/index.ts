// src/services/api/index.ts
/**
 * API Services - Centralized Export Point
 * 
 * This module exports all API endpoint functions for the web-admin application.
 * Each service file corresponds to a specific microservice in the backend.
 * 
 * Supported Services:
 * - Auth: Login, logout, token refresh, device activation
 * - Identity & Organization: Employee/User management, face data, companies
 * - Device Management: Device CRUD operations
 * - Attendance: Check-in, check-out, record retrieval
 * - Workforce: Shifts and schedules management
 * - Analytics & Reporting: Daily reports, summary reports, exports
 * - Signature Upload: User signature management
 * - Account: Personal profile and account settings
 * 
 * Usage:
 * ```
 * import { loginAPI, getCurrentUserAPI } from '@/services/api';
 * import { getEmployeesAPI, createEmployeeAPI } from '@/services/api';
 * import { checkInAPI, checkOutAPI } from '@/services/api';
 * import { getDailyReportAPI, getSummaryReportAPI } from '@/services/api';
 * ```
 */

// ===== Authentication Service =====
export {
  loginAPI,
  logoutAPI,
  refreshTokenAPI,
  getCurrentUserAPI,
  activateDeviceAPI
} from './auth.api';

// ===== Identity & Organization Service (Users) =====
export {
  getEmployeesAPI,
  getEmployeeAPI,
  createEmployeeAPI,
  updateEmployeeAPI,
  deleteEmployeeAPI,
  getEmployeeFaceDataAPI,
  uploadFaceDataAPI,
  deleteFaceDataAPI
} from './employees.api';

// ===== Device Management Service =====
export {
  getDevicesAPI,
  getDeviceAPI,
  createDeviceAPI,
  updateDeviceAPI,
  deleteDeviceAPI,
  syncDeviceAPI
} from './devices.api';

// ===== Attendance Service =====
export {
  checkInAPI,
  checkOutAPI,
  getAttendanceRecordsAPI,
  getAttendanceRecordAPI,
  getMyAttendanceHistoryAPI
} from './attendance.api';

// ===== Workforce Service =====
export {
  getShiftsAPI,
  getShiftAPI,
  createShiftAPI,
  updateShiftAPI,
  deleteShiftAPI,
  getSchedulesAPI,
  getScheduleAPI,
  createScheduleAPI,
  updateScheduleAPI,
  deleteScheduleAPI
} from './shifts.api';

// ===== Analytics & Reporting Service =====
export {
  getDailyReportAPI,
  getSummaryReportAPI,
  exportReportAPI,
  type DailyReport,
  type SummaryReport,
  type ExportReportParams
} from './reports.api';

// ===== Signature Upload Service =====
export {
  uploadSignatureAPI,
  getUserSignatureAPI,
  deleteSignatureAPI,
  type SignatureData
} from './signatures.api';

// ===== Account Management =====
export {
  getAccountProfileAPI,
  updateAccountProfileAPI,
  changeAccountPasswordAPI,
  uploadAccountAvatarAPI
} from './account.api';
