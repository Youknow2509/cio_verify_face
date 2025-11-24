CREATE TABLE IF NOT EXISTS attendance_records (
    company_id UUID,
    year_month TEXT,
    record_time TIMESTAMP,
    employee_id UUID,
    device_id UUID,
    record_type INT,          
    verification_method TEXT,
    verification_score float,
    face_image_url TEXT,
    location_coordinates TEXT,
    metadata MAP<TEXT, TEXT>,
    sync_status TEXT,        
    created_at TIMESTAMP,
    PRIMARY KEY ((company_id, year_month), record_time, employee_id)
) WITH CLUSTERING ORDER BY (record_time DESC);