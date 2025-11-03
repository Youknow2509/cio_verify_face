-- name: CheckUserIsManagementInCompany :one
SELECT e.employee_code
FROM employees e
JOIN users u ON e.employee_id = u.user_id
WHERE e.company_id = $1
  AND e.employee_id = $2
  AND u.role < 2;

-- name: CheckDeviceExistInCompany :one
SELECT d.name
FROM devices d
WHERE d.company_id = $1
  AND d.device_id = $2;

-- name: UpdateDeviceSession :exec
UPDATE devices
SET token = $2,
    updated_at = NOW()
WHERE device_id = $1;

-- name: DisableDevice :exec
UPDATE devices
SET is_active = FALSE,
    updated_at = NOW()
WHERE device_id = $1;

-- name: DeleteDeviceSessionByID :exec
UPDATE devices
SET token = NULL,
    updated_at = NOW()
WHERE device_id = $1;