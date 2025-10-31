CREATE TABLE IF NOT EXISTS attendance_records_by_company_date (
    company_id UUID,
    record_date DATE,
    record_time TIMESTAMP,
    record_id UUID,
    employee_id UUID,
    device_id UUID,
    record_type INT, -- 0: CHECK_IN, 1: CHECK_OUT
    verification_method TEXT,
    verification_score DECIMAL,
    face_image_url TEXT,
    location_coordinates TEXT, -- "lat,lng" format
    metadata MAP<TEXT, TEXT>,
    sync_status TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY ((company_id, record_date), record_time, record_id)
) WITH CLUSTERING ORDER BY (record_time DESC)
AND compression = {'class': 'LZ4Compressor'}
AND compaction = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 1
}
AND gc_grace_seconds = 864000  -- 10 days
AND default_time_to_live = 2592000; -- 30 days TTL

-- Index for employee-based queries
CREATE INDEX IF NOT EXISTS attendance_employee_idx 
ON attendance_records_by_company_date (employee_id);

-- Index for device-based queries
CREATE INDEX IF NOT EXISTS attendance_device_idx 
ON attendance_records_by_company_date (device_id);