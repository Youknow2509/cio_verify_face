// Employee Services
export {
  getEmployees,
  createEmployee,
  updateEmployee,
  deleteEmployee,
  getEmployeeFaceData,
  uploadFaceData,
  deleteFaceData
} from './mock/employees';

// Auth Services  
export {
  login,
  logout,
  refreshToken,
  getCurrentUser
} from './mock/auth';

// Device Services
export {
  getDevices,
  createDevice,
  updateDevice,
  deleteDevice
} from './mock/devices';

// Attendance Services
export {
  getAttendanceRecords,
  getEmployeeAttendance,
  getAttendanceChart
} from './mock/attendance';

// Shift Services
export {
  getShifts,
  createShift,
  updateShift,
  deleteShift
} from './mock/shifts';