-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS companies_secret (
    company_secret_id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    salt VARCHAR(255) NOT NULL,
    secret_hash VARCHAR(512) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Partition and index ...
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS companies_secret;
-- +goose StatementEnd
