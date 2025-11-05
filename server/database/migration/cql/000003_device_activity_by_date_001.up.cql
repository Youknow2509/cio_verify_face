-- Track device heartbeat and activity - live 90 day
CREATE TABLE IF NOT EXISTS device_activity_by_date (
    device_id uuid,
    activity_date date,
    heartbeat_time timestamp,
    device_status text,
    firmware_version text,
    ip_address text,
    last_record_count int,
    error_count int,
    metadata map<text, text>,
    PRIMARY KEY ((device_id, activity_date), heartbeat_time)
) WITH CLUSTERING ORDER BY (heartbeat_time DESC)
AND default_time_to_live = 7776000; 