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