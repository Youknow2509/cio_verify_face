-- =================================================================
-- PASSWORD RESET REQUEST QUERIES
-- =================================================================

-- name: CreatePasswordResetRequest :exec
INSERT INTO password_reset_requests (
    request_id,
    user_id,
    company_id,
    requested_by,
    status,
    meta_data,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPasswordResetRequestByID :one
SELECT 
    request_id,
    user_id,
    company_id,
    requested_by,
    status,
    kafka_message_id,
    kafka_sent_at,
    meta_data,
    created_at,
    updated_at
FROM password_reset_requests
WHERE request_id = $1
LIMIT 1;

-- name: GetRecentRequestsByManagerForUser :many
SELECT 
    request_id,
    user_id,
    company_id,
    requested_by,
    status,
    kafka_message_id,
    kafka_sent_at,
    meta_data,
    created_at,
    updated_at
FROM password_reset_requests
WHERE requested_by = $1 AND user_id = $2 AND created_at >= $3
ORDER BY created_at DESC;

-- name: UpdatePasswordResetStatus :exec
UPDATE password_reset_requests
SET 
    status = $2,
    kafka_message_id = $3,
    kafka_sent_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE request_id = $1;

-- name: CountRequestsByManagerInWindow :one
SELECT COUNT(*) as count
FROM password_reset_requests
WHERE requested_by = $1 AND created_at >= $2;
