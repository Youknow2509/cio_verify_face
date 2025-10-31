/**
 * Services Module
 * 
 * Main entry point for all services:
 * - API services (src/services/api/)
 * - HTTP client (src/services/http.ts)
 * - Error handling (src/services/error-handler.ts)
 * - Mock services (src/services/mock/) - for fallback/testing
 */

// === API Services ===
// These are the main services to use for API calls

export * from '@/services/api';
export { http, setAuthToken, clearAuthToken } from '@/services/http';
export type { RequestConfig, HttpError } from '@/services/http';

// === Error Handling ===
export {
  handleApiError,
  isAuthError,
  isValidationError,
  isNotFoundError,
  isServerError,
  getUserFriendlyMessage,
  type ErrorInfo
} from '@/services/error-handler';

// === HTTP Interceptor ===
export {
  configureHttpInterceptor,
  logRequest,
  logResponse,
  logError,
  isRetryableError,
  getRetryDelay,
  sleep,
  type InterceptorConfig
} from '@/services/http-interceptor';

// === Wrapper Functions (for backward compatibility) ===
// Re-export API functions with shorter names
export {
  getEmployeesAPI as getEmployees,
  getEmployeeAPI as getEmployee,
  createEmployeeAPI as createEmployee,
  updateEmployeeAPI as updateEmployee,
  deleteEmployeeAPI as deleteEmployee,
  getEmployeeFaceDataAPI as getEmployeeFaceData,
  uploadFaceDataAPI as uploadFaceData,
  deleteFaceDataAPI as deleteFaceData,
} from '@/services/api/employees.api';

export {
  getDevicesAPI as getDevices,
  getDeviceAPI as getDevice,
  createDeviceAPI as createDevice,
  updateDeviceAPI as updateDevice,
  deleteDeviceAPI as deleteDevice,
  syncDeviceAPI as syncDevice,
} from '@/services/api/devices.api';

export {
  getShiftsAPI as getShifts,
  getShiftAPI as getShift,
  createShiftAPI as createShift,
  updateShiftAPI as updateShift,
  deleteShiftAPI as deleteShift,
  getSchedulesAPI as getSchedules,
  getScheduleAPI as getSchedule,
  createScheduleAPI as createSchedule,
  updateScheduleAPI as updateSchedule,
  deleteScheduleAPI as deleteSchedule,
} from '@/services/api/shifts.api';

export {
  getAttendanceRecordsAPI as getAttendanceRecords,
  getAttendanceRecordAPI as getAttendanceRecord,
  getMyAttendanceHistoryAPI as getMyAttendanceHistory,
  checkInAPI as checkIn,
  checkOutAPI as checkOut,
} from '@/services/api/attendance.api';

export {
  getDailyReportAPI as getDailyReport,
  getSummaryReportAPI as getSummaryReport,
  exportReportAPI as exportReport,
} from '@/services/api/reports.api';

// === Mock Services (Legacy) ===
// Import from these only for testing/fallback scenarios
// Note: Mock services are deprecated, use API services above
export * as mockServices from './mock';