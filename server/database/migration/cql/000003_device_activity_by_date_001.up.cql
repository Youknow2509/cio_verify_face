-- Track device heartbeat and activity
CREATE TABLE IF NOT EXISTS device_activity_by_date (
    device_id UUID,
    activity_date DATE,
    heartbeat_time TIMESTAMP,
    device_status TEXT,
    firmware_version TEXT,
    ip_address TEXT,
    last_record_count INT,
    error_count INT,
    metadata MAP<TEXT, TEXT>,
    PRIMARY KEY ((device_id, activity_date), heartbeat_time)
) WITH CLUSTERING ORDER BY (heartbeat_time DESC)
AND compression = {'class': 'LZ4Compressor'}
AND default_time_to_live = 7776000; -- 90 days TTL