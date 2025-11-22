CREATE TABLE IF NOT EXISTS attendance_records_no_shift (
    company_id UUID,
    year_month TEXT,
    record_time TIMESTAMP,
    employee_id UUID,
    device_id UUID,
    verification_method TEXT,
    verification_score float,
    face_image_url TEXT,
    location_coordinates TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY ((company_id, year_month), record_time)
) WITH CLUSTERING ORDER BY (record_time DESC);