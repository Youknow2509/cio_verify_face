CREATE TABLE IF NOT EXISTS user_verifications (
    company_id UUID,
    user_id UUID,
    timestamp TIMESTAMP,
    verification_id UUID,
    verified BOOLEAN,
    similarity_score DOUBLE,
    device_id TEXT,
    PRIMARY KEY ((company_id, user_id), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);
