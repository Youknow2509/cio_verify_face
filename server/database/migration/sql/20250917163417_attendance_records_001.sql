-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS attendance_records (
    record_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(employee_id) ON DELETE CASCADE,
    device_id UUID NOT NULL REFERENCES devices(device_id) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    record_type INTEGER NOT NULL CHECK (record_type IN (0, 1)), -- 0: CHECK_IN, 1: CHECK_OUT
    verification_method VARCHAR(20) DEFAULT 'FACE' CHECK (verification_method IN ('FACE', 'MANUAL', 'CARD')),
    verification_score DECIMAL(5,3), -- Face match confidence score
    face_image_url VARCHAR(500),
    location_coordinates POINT, -- GPS coordinates if available
    metadata JSONB DEFAULT '{}', -- Additional data: temperature, mask detection, etc.
    sync_status int2 DEFAULT 1 CHECK (sync_status IN (0, 1, 2)), -- 0: SYNCED, 1: PENDING, 2: FAILED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attendance_records;
-- +goose StatementEnd
