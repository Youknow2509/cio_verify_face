CREATE TABLE IF NOT EXISTS enrollment_states (
    profile_id UUID,
    company_id UUID,
    user_id UUID,
    device_id TEXT,
    status TEXT,
    quality_score DOUBLE,
    timestamp TIMESTAMP,
    metadata MAP<TEXT, TEXT>,
    image_path TEXT,
    PRIMARY KEY ((profile_id, company_id), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);
