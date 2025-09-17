-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS company_settings (
    setting_id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    setting_key VARCHAR(100) NOT NULL,
    setting_value TEXT,
    setting_type int2 DEFAULT 0 CHECK (setting_type IN (0, 1, 2, 3)), -- 0: STRING, 1: INTEGER, 2: BOOLEAN, 3: JSON
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(company_id, setting_key)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS company_settings;
-- +goose StatementEnd
