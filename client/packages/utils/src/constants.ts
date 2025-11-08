export const ATTENDANCE_STATUS = {
  ON_TIME: 'on_time',
  LATE: 'late',
  EARLY: 'early',
} as const;

export const DEVICE_STATUS = {
  ONLINE: 'online',
  OFFLINE: 'offline',
  ERROR: 'error',
} as const;

export const USER_ROLE = {
  SYSTEM_ADMIN: 'system_admin',
  COMPANY_ADMIN: 'company_admin',
  EMPLOYEE: 'employee',
} as const;

export const FACE_QUALITY = {
  GOOD: 'good',
  AVERAGE: 'average',
  POOR: 'poor',
} as const;
