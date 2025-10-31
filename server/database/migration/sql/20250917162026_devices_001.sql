-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS devices (
    device_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    location_id UUID, -- For future location management
    name VARCHAR(255) NOT NULL,
    address TEXT,
    device_type int2 DEFAULT 0 CHECK (device_type IN (0, 1, 2, 3)), -- 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
    serial_number VARCHAR(100),
    mac_address VARCHAR(17),
    ip_address INET,
    firmware_version VARCHAR(20),
    status int2 DEFAULT 1 CHECK (status IN (0, 1, 2, 3)), -- 0: OFFLINE, 1: ONLINE, 2: MAINTENANCE, 3: ERROR
    token VARCHAR(512) NOT NULL, -- Device authentication token
    last_heartbeat TIMESTAMP WITH TIME ZONE,
    settings JSONB DEFAULT '{}', -- Device-specific settings
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS devices;
-- +goose StatementEnd
