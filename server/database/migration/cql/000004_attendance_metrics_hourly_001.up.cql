-- Hourly attendance metrics for analytics
CREATE TABLE IF NOT EXISTS attendance_metrics_hourly (
    company_id UUID,
    metric_date DATE,
    metric_hour INT,
    total_checkins INT,
    total_checkouts INT,
    unique_employees INT,
    active_devices INT,
    avg_verification_score DECIMAL,
    peak_concurrent_users INT,
    created_at TIMESTAMP,
    PRIMARY KEY ((company_id, metric_date), metric_hour)
) WITH CLUSTERING ORDER BY (metric_hour ASC)
AND compression = {'class': 'LZ4Compressor'}
AND default_time_to_live = 2592000; -- 30 days TTL
