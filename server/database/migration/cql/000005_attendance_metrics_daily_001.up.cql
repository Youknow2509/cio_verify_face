
-- Daily rollup metrics
CREATE TABLE IF NOT EXISTS attendance_metrics_daily (
    company_id UUID,
    metric_date DATE,
    total_attendance_records INT,
    unique_employees_count INT,
    present_count INT,
    late_count INT,
    absent_count INT,
    avg_work_hours DECIMAL,
    total_overtime_minutes INT,
    attendance_rate DECIMAL,
    punctuality_rate DECIMAL,
    created_at TIMESTAMP,
    PRIMARY KEY (company_id, metric_date)
) WITH CLUSTERING ORDER BY (metric_date DESC)
AND compression = {'class': 'LZ4Compressor'}
AND default_time_to_live = 31536000; -- 1 year TTL