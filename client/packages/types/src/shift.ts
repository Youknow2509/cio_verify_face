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
