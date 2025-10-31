-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS work_shifts (
    shift_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_duration_minutes INTEGER DEFAULT 0,
    grace_period_minutes INTEGER DEFAULT 15, -- Late arrival tolerance
    early_departure_minutes INTEGER DEFAULT 15, -- Early leave tolerance
    work_days INTEGER[] DEFAULT ARRAY[1,2,3,4,5], -- 1=Monday, 7=Sunday
    is_flexible BOOLEAN DEFAULT FALSE,
    overtime_after_minutes INTEGER DEFAULT 480, -- 8 hours = 480 minutes
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS work_shifts;
-- +goose StatementEnd
