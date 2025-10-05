-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS daily_attendance_summary (
    summary_id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    employee_id UUID NOT NULL REFERENCES employees(employee_id) ON DELETE CASCADE,
    shift_id UUID REFERENCES work_shifts(shift_id),
    work_date DATE NOT NULL,
    scheduled_in TIME,
    scheduled_out TIME,
    actual_check_in TIMESTAMP WITH TIME ZONE,
    actual_check_out TIMESTAMP WITH TIME ZONE,
    total_work_minutes INTEGER DEFAULT 0,
    break_minutes INTEGER DEFAULT 0,
    overtime_minutes INTEGER DEFAULT 0,
    late_minutes INTEGER DEFAULT 0,
    early_leave_minutes INTEGER DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 3 CHECK (status IN (0, 1, 2, 3)), -- 0: PRESENT, 1: LATE, 2: EARLY_LEAVE, 3: ABSENT
    attendance_percentage DECIMAL(5,2) DEFAULT 0.00,
    notes TEXT,
    approved_by UUID REFERENCES Users(user_id),
    approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(employee_id, work_date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS daily_attendance_summary;
-- +goose StatementEnd
