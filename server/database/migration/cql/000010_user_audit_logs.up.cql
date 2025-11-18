CREATE TABLE IF NOT EXISTS user_audit_logs (
    company_id UUID,
    user_id UUID,
    timestamp TIMESTAMP,
    log_id UUID,
    operation TEXT,
    status TEXT,
    profile_id UUID,
    device_id TEXT,
    similarity_score DOUBLE,
    liveness_score DOUBLE,
    quality_score DOUBLE,
    error_message TEXT,
    PRIMARY KEY ((company_id, user_id), timestamp, log_id)
) WITH CLUSTERING ORDER BY (timestamp DESC, log_id ASC);
