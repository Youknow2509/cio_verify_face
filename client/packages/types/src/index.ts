// Shared TypeScript types for Face Attendance SaaS

// User & Authentication Types
export interface User {
    id: string;
    email: string;
    full_name: string;
    employee_code?: string;
    phone?: string;
    avatar_url?: string;
    role: 'system_admin' | 'company_admin' | 'employee';
    company_id?: string;
    status: 'active' | 'inactive';
    created_at: string;
    updated_at: string;
}

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    code: Int16Array | string;
    message: string;
    data: {
        access_token: string;
        refresh_token: string;
    };
}

export interface AuthState {
    user: User | null;
    accessToken: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;
}

// Company Types
export interface Company {
    id: string;
    name: string;
    email: string;
    phone?: string;
    address?: string;
    website?: string;
    logo_url?: string;
    status: 'active' | 'suspended' | 'expired';
    max_employees: number;
    subscription_plan: string;
    subscription_start: string;
    subscription_end: string;
    created_at: string;
    updated_at: string;
}

// Employee Types
export interface Employee extends User {
    position?: string;
    department?: string;
    hire_date?: string;
    contract_type?: string;
    default_shift_id?: string;
    face_data_count: number;
}

export interface FaceData {
    id: string;
    user_id: string;
    image_url: string;
    quality: 'good' | 'average' | 'poor';
    uploaded_at: string;
}

// Device Types
export interface Device {
    id: string;
    company_id: string;
    name: string;
    serial_number: string;
    location: string;
    description?: string;
    ip_address?: string;
    mac_address?: string;
    status: 'online' | 'offline' | 'error';
    last_online_at?: string;
    settings: DeviceSettings;
    created_at: string;
    updated_at: string;
    token: string;
}

// Device editing name
export interface DeviceEditNameForm {
    device_id: string;
    device_name: string;
}

// Device editing location
export interface DeviceEditLocationForm {
    device_id: string;
    location_id?: string;
    address: string;
}

// Device creation
export interface DeviceCreateRequest {
    address: string;
    device_name: string;
    device_type: DeviceType; // numeric enum provided by backend
    mac_address: string;
    serial_number: string;
}

export enum DeviceType {
    FACE_TERMINAL = 0,
    MOBILE_APP = 1,
    WEB_CAMERA = 2,
    IOT_SENSOR = 3,
}

export interface DeviceSettings {
    allow_check_in: boolean;
    allow_check_out: boolean;
    timeout: number;
    recognition_threshold: number;
    sound_enabled: boolean;
}

// Shift Types
export interface Shift {
    shift_id: string;
    company_id: string;
    name: string;
    description?: string;
    start_time: string; // HH:MM format
    end_time: string; // HH:MM format
    break_duration_minutes: number;
    grace_period_minutes: number; // Late arrival tolerance
    early_departure_minutes: number; // Early leave tolerance
    work_days: number[]; // 1=Monday, 7=Sunday
    is_flexible: boolean;
    overtime_after_minutes: number; // Default 480 (8 hours)
    is_active: boolean;
    created_at: string;
    updated_at: string;

    // Optional fields for UI
    employee_count?: number;
}

export interface EmployeeShift {
    employee_id: string;
    shift_id: string;
    effective_from: string; // Date format
    effective_to?: string; // Date format
    is_active: boolean;
    created_at: string;
    // Optional fields for UI
    employee_name?: string;
    shift_name?: string;
}

export interface ShiftFormData {
    name: string;
    description: string;
    start_time: string;
    end_time: string;
    break_duration_minutes: number;
    grace_period_minutes: number;
    early_departure_minutes: number;
    work_days: number[];
    is_flexible: boolean;
    overtime_after_minutes: number;
    is_active: boolean;
}

export interface Schedule {
    id: string;
    employee_id: string;
    shift_id: string;
    start_date: string;
    end_date?: string;
    schedule_type: 'permanent' | 'temporary';
    created_at: string;
}

// Attendance Types
export interface AttendanceRecord {
    id: string;
    employee_id: string;
    employee_name: string;
    employee_code: string;
    device_id: string;
    device_name: string;
    shift_id?: string;
    shift_name?: string;
    check_in_time?: string;
    check_out_time?: string;
    check_in_status?: 'on_time' | 'late' | 'early';
    check_out_status?: 'on_time' | 'early' | 'late';
    total_hours?: number;
    date: string;
    notes?: string;
    created_at: string;
    updated_at: string;
}

export interface AttendanceCheckRequest {
    device_id: string;
    face_image: string; // base64
    check_type: 'check_in' | 'check_out';
    timestamp: string;
}

export interface AttendanceCheckResponse {
    success: boolean;
    employee?: {
        id: string;
        name: string;
        code: string;
        avatar_url?: string;
    };
    check_type: 'check_in' | 'check_out';
    timestamp: string;
    status: 'on_time' | 'late' | 'early';
    message: string;
    error?: string;
}

// Report Types
export interface DailyReport {
    date: string;
    total_employees: number;
    present_count: number;
    absent_count: number;
    late_count: number;
    early_leave_count: number;
    on_time_rate: number;
    records: AttendanceRecord[];
}

export interface SummaryReportFilter {
    start_date: string;
    end_date: string;
    employee_ids?: string[];
    department?: string;
    group_by: 'employee' | 'department' | 'device';
}

export interface SummaryReportData {
    employee_id?: string;
    employee_name?: string;
    department?: string;
    total_days: number;
    total_hours: number;
    late_count: number;
    early_leave_count: number;
    compliance_rate: number;
}

// Dashboard Types
export interface DashboardStats {
    total_employees: number;
    attendance_today: number;
    late_rate_this_month: number;
    active_devices: number;
}

export interface DashboardChartData {
    date: string;
    check_ins: number;
    check_outs: number;
}

// Settings Types
export interface CompanySettings {
    company_id: string;
    attendance_settings: {
        valid_time_before_shift: number;
        valid_time_after_shift: number;
        allow_offline_attendance: boolean;
        recognition_threshold: number;
        timeout: number;
    };
    notification_settings: {
        daily_email_report: boolean;
        device_offline_alert: boolean;
        abnormal_attendance_alert: boolean;
    };
}

// API Response Types
export interface ApiResponse<T> {
    success: boolean;
    data?: T;
    message?: string;
    error?: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    per_page: number;
    total_pages: number;
}

// WebSocket Event Types
export interface WebSocketEvent {
    type: 'attendance_result' | 'device_status' | 'admin_alert';
    data: any;
    timestamp: string;
}

// System Admin Types
export interface SystemStats {
    total_companies: number;
    total_employees: number;
    total_devices: number;
    total_attendance_today: number;
}

export interface AuditLog {
    id: string;
    user_id: string;
    user_name: string;
    company_id?: string;
    action: string;
    resource: string;
    details: Record<string, any>;
    ip_address: string;
    timestamp: string;
}

// Export all shift types
export * from './shift';
// Face profile types
export * from './faceProfile';
