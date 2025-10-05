// src/types/index.ts

export type UserRole = 'CompanyAdmin' | 'Manager' | 'Staff';

export interface Employee {
  id: string;
  code: string;
  name: string;
  email: string;
  department?: string;
  position?: string;
  active: boolean;
  faceCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface FaceData {
  id: string;
  employeeId: string;
  imageUrl: string;
  fileName: string;
  createdAt: string;
}

export interface Device {
  id: string;
  name: string;
  location?: string;
  status: 'online' | 'offline';
  lastSyncAt?: string;
  model?: string;
  ipAddress?: string;
  createdAt: string;
}

export interface AttendanceRecord {
  id: string;
  employeeId: string;
  employeeName: string;
  date: string;
  checkIn?: string;
  checkOut?: string;
  isLate?: boolean;
  shiftId?: string;
  totalHours?: number;
  deviceId?: string;
}

export interface Shift {
  id: string;
  name: string;
  start: string;
  end: string;
  description?: string;
  active: boolean;
  createdAt: string;
}

export interface ReportRow {
  employeeId: string;
  employeeName: string;
  date: string;
  totalHours: number;
  lateMinutes: number;
  department?: string;
  checkIn?: string;
  checkOut?: string;
}

export interface Company {
  id: string;
  name: string;
  logo?: string;
  timezone: string;
  dateFormat: string;
  locale: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  companyId: string;
  active: boolean;
}

export interface ApiResponse<T> {
  data: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface FilterOptions {
  search?: string;
  department?: string;
  status?: string;
  startDate?: string;
  endDate?: string;
  page?: number;
  limit?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface EmployeeFilter {
  page?: number;
  limit?: number;
  search?: string;
  department?: string;
  active?: boolean;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface TableColumn<T = any> {
  key: keyof T | string;
  header: string;
  sortable?: boolean;
  width?: string;
  align?: 'left' | 'center' | 'right';
  render?: (value: any, record: T, index: number) => React.ReactNode;
}

export interface DashboardStats {
  totalEmployees: number;
  todayCheckIns: number;
  lateArrivals: number;
  devicesOnline: number;
  attendanceRate: number;
}

export interface ChartData {
  date: string;
  checkIns: number;
  checkOuts: number;
  lateArrivals: number;
}

export interface RecentActivity {
  id: string;
  type: 'check_in' | 'check_out' | 'device_sync' | 'employee_added';
  message: string;
  timestamp: string;
  employeeName?: string;
  deviceName?: string;
}