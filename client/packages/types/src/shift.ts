// Work Shift Types
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
    employee_shift_id: string;
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

// API Request/Response Types
export interface ResponseData<T = any> {
    code: number;
    data?: T;
    message?: string;
}

export interface ErrResponseData {
    code: number;
    error: string;
    detail?: any;
}

// Shift API Request DTOs
export interface CreateShiftReq {
    company_id: string;
    name: string;
    description?: string;
    start_time: number; // Unix timestamp in seconds
    end_time: number; // Unix timestamp in seconds
    break_duration_minutes: number;
    grace_period_minutes: number;
    early_departure_minutes: number;
    work_days: number[]; // 1=Monday, 7=Sunday
}

export interface EditShiftReq extends CreateShiftReq {
    shift_id: string;
}

export interface ChangeStatusShiftReq {
    company_id: string;
    shift_id: string;
    status: 0 | 1; // 0=inactive, 1=active
}

// Employee Shift API Request DTOs
export interface AddShiftEmployeeReq {
    employee_id: string;
    shift_id: string;
    effective_from: number; // Unix timestamp
    effective_to: number; // Unix timestamp
}

export interface AddShiftEmployeeListReq {
    company_id: string;
    employee_ids: string[];
    shift_id: string;
    effective_from: number; // Unix timestamp
    effective_to: number; // Unix timestamp
}

export interface ShiftEmployeeEffectiveDateReq {
    user_id?: string;
    effective_from: number; // Unix timestamp
    effective_to: number; // Unix timestamp
    page: number;
    size: number;
}

export interface ShiftEmployeeEditEffectiveDateReq {
    shift_user_id: string;
    new_effective_from: number; // Unix timestamp
    new_effective_to: number; // Unix timestamp
}

export interface EnableShiftForUserReq {
    shift_user_id: string;
}

export interface DisableShiftForUserReq {
    shift_user_id: string;
}

// ==================== New Shift Employee Query DTOs ====================
// Assumptions based on backend swagger for /v1/employee/shift/in & /v1/employee/shift/not_in
// If backend adds more filters (e.g. department), extend these interfaces accordingly.
export interface GetInfoEmployeeInShiftReq {
    shift_id: string;
    page: number;
    size: number;
}

export interface GetInfoEmployeeNotInShiftReq {
    shift_id: string;
    page: number;
    size: number;
}


export interface DeleteShiftEmployeeReq {
    employee_ids: string[];
    shift_id: string;
}