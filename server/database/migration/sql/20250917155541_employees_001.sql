-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employees (
    employee_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    employee_code VARCHAR(50) NOT NULL,
    department VARCHAR(100),
    position VARCHAR(100),
    hire_date DATE,
    salary DECIMAL(12,2),
    status int2 DEFAULT 0 NOT NULL CHECK (status IN (0, 1, 2)), -- 0: active, 1: inactive, 2: on leave
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
-- +goose StatementEnd

