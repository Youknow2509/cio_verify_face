export interface Company {
  company_id: string;
  name: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  status: number; // 0: Inactive, 1: Active, 2: Suspended
  subscription_plan: number; // 0: Basic, 1: Premium, 2: Enterprise
  subscription_start_date?: string;
  subscription_end_date?: string;
  max_employees: number;
  max_devices: number;
  created_at: string;
  updated_at: string;
}

export interface User {
  user_id: string;
  email: string;
  phone: string;
  salt: string;
  password_hash: string;
  full_name: string;
  avatar_url?: string;
  role: number; // 0: SYSTEM_ADMIN, 1: COMPANY_ADMIN, 2: EMPLOYEE
  status: number; // 0: ACTIVE, 1: INACTIVE, 2: SUSPENDED
  last_login?: string;
  is_locked: boolean;
  lock_expires_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Employee {
  employee_id: string;
  company_id: string;
  employee_code: string;
  department?: string;
  position?: string;
  hire_date?: string;
  salary?: number;
  status: number; // 0: active, 1: inactive, 2: on leave
  created_at: string;
  updated_at: string;
}

export interface FaceData {
  fid: string;
  user_id: string;
  image_url: string;
  face_encoding?: string;
  quality_score?: number;
  created_at: string;
  updated_at: string;
}

export interface CreateCompanyRequest {
  name: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  status?: number;
  subscription_plan?: number;
  subscription_start_date?: string;
  subscription_end_date?: string;
  max_employees?: number;
  max_devices?: number;
}

export interface UpdateCompanyRequest {
  name?: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  status?: number;
  subscription_plan?: number;
  subscription_start_date?: string;
  subscription_end_date?: string;
  max_employees?: number;
  max_devices?: number;
}

export interface CreateUserRequest {
  email: string;
  phone: string;
  password: string;
  full_name: string;
  avatar_url?: string;
  role: number;
  company_id?: string;
  employee_code?: string;
  department?: string;
  position?: string;
  hire_date?: string;
  salary?: number;
}

export interface UpdateUserRequest {
  email?: string;
  phone?: string;
  full_name?: string;
  avatar_url?: string;
  status?: number;
  department?: string;
  position?: string;
  salary?: number;
}

export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
  statusCode: number;
}
