CREATE TABLE IF NOT EXISTS attendance_records_by_user (
    company_id UUID,
    employee_id UUID,
    year_month TEXT,        
    record_time TIMESTAMP,
    device_id UUID,
    record_type INT,
    verification_method TEXT,
    verification_score float,
    face_image_url TEXT,
    location_coordinates TEXT,
    metadata MAP<TEXT, TEXT>,
    sync_status TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY ((company_id, employee_id, year_month), record_time)
) WITH CLUSTERING ORDER BY (record_time DESC);