-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employee_shifts (
    employee_shift_id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    employee_id UUID NOT NULL REFERENCES employees(employee_id) ON DELETE CASCADE,
    shift_id UUID NOT NULL REFERENCES work_shifts(shift_id) ON DELETE CASCADE,
    effective_from DATE NOT NULL,
    effective_to DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, shift_id, effective_from)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employee_shifts;
-- +goose StatementEnd
