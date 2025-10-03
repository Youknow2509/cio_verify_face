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