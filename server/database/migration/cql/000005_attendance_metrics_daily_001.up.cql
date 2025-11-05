
-- Daily rollup metrics - 1 year TTL
CREATE TABLE IF NOT EXISTS attendance_metrics_daily (
    company_id UUID,
    metric_date DATE,
    total_attendance_records INT,
    unique_employees_count INT,
    present_count INT,
    late_count INT,
    absent_count INT,
    avg_work_hours float,
    total_overtime_minutes INT,
    attendance_rate float,
    punctuality_rate float,
    created_at TIMESTAMP,
    PRIMARY KEY (company_id, metric_date)
) WITH CLUSTERING ORDER BY (metric_date DESC)
AND default_time_to_live = 31536000;