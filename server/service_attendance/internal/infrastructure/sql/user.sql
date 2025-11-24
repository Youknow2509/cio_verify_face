-- name: GetListTimeShiftEmployee :many
SELECT 
    ws.shift_id,
    ws.start_time,
    ws.end_time,
    ws.grace_period_minutes,
    ws.early_departure_minutes,
    ws.work_days,
    es.effective_from,
    es.effective_to
FROM 
    employee_shifts es
JOIN 
    work_shifts ws ON es.shift_id = ws.shift_id
WHERE 
    ws.company_id = $1
    AND es.employee_id = $2
    AND es.is_active = TRUE
    AND ws.is_active = TRUE;

-- name: GetCompanyIdUser :one
SELECT company_id
FROM employees
WHERE employee_id = $1
LIMIT 1;

-- name: GetUserInfoWithID :one
SELECT 
    email,
    phone,
    full_name,
    avatar_url
FROM users
WHERE user_id = $1
LIMIT 1;

-- name: GetUserBaseWithMail :one
SELECT user_id, email, salt, password_hash, role, is_locked 
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreateUserSession :exec
INSERT INTO user_sessions (
    session_id, 
    user_id, 
    refresh_token,
    ip_address,
    user_agent,
    created_at,
    expires_at
) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6);

-- name: DeleteUserSessionByID :exec
DELETE FROM user_sessions
WHERE session_id = $1;

-- name: GetUserSessionByID :one
SELECT 
    session_id, 
    user_id, 
    refresh_token,
    ip_address,
    user_agent,
    created_at,
    expires_at
FROM user_sessions
WHERE session_id = $1
LIMIT 1;

-- name: UpdateUserSession :exec
UPDATE user_sessions
SET
    refresh_token = $2,
    expires_at = $3
WHERE session_id = $1;