
-- name: BlockDeviceToken :exec
UPDATE devices
SET
    token = NULL,
    status = 0,
    updated_at = NOW()
WHERE device_id = $1;

-- name: CreateDeviceToken :exec
UPDATE devices
SET
    token = $2,
    status = 1,
    updated_at = NOW()
WHERE device_id = $1;

-- name: DeviceExists :one
SELECT device_id
FROM devices
WHERE device_id = $1;

-- name: CheckTokenDevice :one
SELECT status
FROM devices
WHERE device_id = $1 AND token = $2;