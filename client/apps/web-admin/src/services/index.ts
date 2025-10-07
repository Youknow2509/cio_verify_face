// Employee Services
export {
  getEmployees,
  createEmployee,
  updateEmployee,
  deleteEmployee,
  getEmployeeFaceData,
  uploadFaceData,
  deleteFaceData
} from '@/services/mock/employees';

// Auth Services  
export {
  login,
  logout,
  refreshToken,
  getCurrentUser
} from '@/services/mock/auth';

// Device Services
export {
  getDevices,
  createDevice,
  updateDevice,
  deleteDevice,
  syncDevice
} from '@/services/mock/devices';

// Attendance Services
export {
  getAttendanceRecords,
  getEmployeeAttendance,
  getAttendanceChart
} from '@/services/mock/attendance';

// Shift Services
export {
  getShifts,
  createShift,
  updateShift,
  deleteShift
} from '@/services/mock/shifts';