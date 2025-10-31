-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS companies (
    company_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    website VARCHAR(255),
    status int2 DEFAULT 0 CHECK (status IN (0, 1, 2)), -- 0: Inactive, 1: Active, 2: Suspended
    subscription_plan int2 DEFAULT 0 CHECK (subscription_plan IN (0, 1, 2)), -- 0: Basic, 1: Premium, 2: Enterprise
    subscription_start_date DATE,
    subscription_end_date DATE,
    max_employees INTEGER DEFAULT 100,
    max_devices INTEGER DEFAULT 10,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Create index and partitioning ...
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Companies;
-- +goose StatementEnd
