CREATE TABLE IF NOT EXISTS authentication_states (
    verification_id UUID,
    company_id UUID,
    user_id UUID,
    profile_id UUID,
    device_id TEXT,
    status TEXT,
    verified BOOLEAN,
    similarity_score DOUBLE,
    liveness_score DOUBLE,
    timestamp TIMESTAMP,
    metadata MAP<TEXT, TEXT>,
    image_path TEXT,
    PRIMARY KEY ((verification_id, company_id), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC)
AND default_time_to_live = 7776000;