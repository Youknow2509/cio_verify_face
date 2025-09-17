-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS system_settings (
    setting_id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT,
    setting_type int2 DEFAULT 0 CHECK (setting_type IN (0, 1, 2, 3)), -- 0: STRING, 1: INTEGER, 2: BOOLEAN, 3: JSON
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE, -- System settings cannot be modified via UI
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS system_settings;
-- +goose StatementEnd
