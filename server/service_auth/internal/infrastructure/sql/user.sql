-- name: GetUserBaseWithMail :one
SELECT user_id, email, salt, password_hash, role, is_locked 
FROM users
WHERE email = $1
LIMIT 1;