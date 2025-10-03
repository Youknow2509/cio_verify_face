-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS attendance_exceptions (
    exception_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    summary_id UUID NOT NULL REFERENCES daily_attendance_summary(summary_id) ON DELETE CASCADE,
    exception_type int2 NOT NULL CHECK (exception_type IN (0, 1, 2, 3)), -- 0: LATE_EXCUSE, 1: EARLY_LEAVE_EXCUSE, 2: OVERTIME_APPROVAL, 3: MANUAL_ADJUSTMENT
    reason TEXT NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(user_id),
    approved_by UUID REFERENCES users(user_id),
    status int2 DEFAULT 0 CHECK (status IN (0, 1, 2)), -- 0: PENDING, 1: APPROVED, 2: REJECTED
    adjustment_minutes INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attendance_exceptions;
-- +goose StatementEnd
