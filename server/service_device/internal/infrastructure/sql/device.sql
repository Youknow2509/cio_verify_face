-- name: CheckDeviceExist :one
SELECT EXISTS (
    SELECT 1
    FROM devices
    WHERE device_id = $1
) AS exist;

-- name: UpdateDeviceInfo :exec
UPDATE devices
SET serial_number = $2,
    mac_address = $3,
    firmware_version = $4,
    updated_at = NOW()
WHERE device_id = $1;

-- name: UpdateDeviceLocation :exec
UPDATE devices
SET location_id = $2, address = $3, updated_at = NOW()
WHERE device_id = $1;

-- name: UpdateDeviceName :exec
UPDATE devices
SET name = $2, updated_at = NOW()
WHERE device_id = $1;

-- name: EnableDevice :exec
UPDATE devices
SET status = 1, updated_at = NOW()
WHERE device_id = $1;

-- name: DisableDevice :exec
UPDATE devices
SET status = 0, updated_at = NOW()
WHERE device_id = $1;

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE device_id = $1;

-- name: GetListDeviceInCompany :many
SELECT
    device_id,
    company_id,
    name,
    address,
    serial_number,
    mac_address,
    status,
    created_at,
    updated_at
FROM devices
WHERE company_id = $1
LIMIT $2 OFFSET $3;

-- name: GetDeviceInfo :one
SELECT
    company_id,
    name,
    address,
    serial_number,
    mac_address,
    ip_address,
    firmware_version,
    last_heartbeat,
    settings,
    created_at,
    updated_at
FROM devices
WHERE device_id = $1
LIMIT 1;

-- name: GetDeviceInfoBase :one
SELECT
    company_id,
    name,
    address,
    serial_number,
    mac_address,
    status,
    created_at,
    updated_at
FROM devices
WHERE device_id = $1
LIMIT 1;

-- name: CreateNewDevice :exec
INSERT INTO devices (
    device_id,
    company_id ,
    name,
    address,
    serial_number,
    mac_address,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), NOW()
);