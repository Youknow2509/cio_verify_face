CREATE TABLE IF NOT EXISTS audit_logs (
    log_id UUID,
    timestamp TIMESTAMP,
    company_id UUID,
    operation TEXT,
    status TEXT,
    profile_id UUID,
    user_id UUID,
    device_id TEXT,
    similarity_score DOUBLE,
    liveness_score DOUBLE,
    quality_score DOUBLE,
    metadata MAP<TEXT, TEXT>,
    error_message TEXT,
    PRIMARY KEY ((company_id, operation), timestamp, log_id)
) WITH CLUSTERING ORDER BY (timestamp DESC, log_id ASC);
