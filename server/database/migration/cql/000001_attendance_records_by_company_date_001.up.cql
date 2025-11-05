CREATE TABLE IF NOT EXISTS attendance_records_by_company_date (
    company_id UUID,
    record_time TIMESTAMP,
    employee_id UUID,
    device_id UUID,
    record_type INT, -- 0: CHECK_IN, 1: CHECK_OUT
    verification_method TEXT,
    verification_score float,
    face_image_url TEXT,
    location_coordinates TEXT, -- "lat,lng" format
    metadata MAP<TEXT, TEXT>,
    sync_status TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY ((company_id, device_id), record_time)
) WITH CLUSTERING ORDER BY (record_time DESC)
AND default_time_to_live = 2592000; -- 30 days TTL
