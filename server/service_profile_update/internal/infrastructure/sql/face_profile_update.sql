-- =================================================================
-- FACE PROFILE UPDATE REQUEST QUERIES
-- =================================================================

-- name: CreateFaceProfileUpdateRequest :exec
INSERT INTO face_profile_update_requests (
    request_id,
    user_id,
    company_id,
    status,
    request_month,
    request_count_in_month,
    reason,
    meta_data,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetFaceProfileUpdateRequestByID :one
SELECT 
    request_id,
    user_id,
    company_id,
    status,
    request_month,
    request_count_in_month,
    update_token,
    update_link_expires_at,
    approved_by,
    approved_at,
    rejection_reason,
    reason,
    meta_data,
    created_at,
    updated_at
FROM face_profile_update_requests
WHERE request_id = $1 AND company_id = $2
LIMIT 1;

-- name: GetFaceProfileUpdateRequestByToken :one
SELECT 
    request_id,
    user_id,
    company_id,
    status,
    request_month,
    request_count_in_month,
    update_token,
    update_link_expires_at,
    approved_by,
    approved_at,
    rejection_reason,
    reason,
    meta_data,
    created_at,
    updated_at
FROM face_profile_update_requests
WHERE update_token = $1
LIMIT 1;

-- name: GetPendingRequestsByCompany :many
SELECT 
    request_id,
    user_id,
    company_id,
    status,
    request_month,
    request_count_in_month,
    update_token,
    update_link_expires_at,
    approved_by,
    approved_at,
    rejection_reason,
    reason,
    meta_data,
    created_at,
    updated_at
FROM face_profile_update_requests
WHERE company_id = $1 AND status = 0
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRequestsByUserAndMonth :many
SELECT 
    request_id,
    user_id,
    company_id,
    status,
    request_month,
    request_count_in_month,
    update_token,
    update_link_expires_at,
    approved_by,
    approved_at,
    rejection_reason,
    reason,
    meta_data,
    created_at,
    updated_at
FROM face_profile_update_requests
WHERE user_id = $1 AND request_month = $2
ORDER BY created_at DESC;

-- name: CountRequestsByUserInMonth :one
SELECT COUNT(*) as count
FROM face_profile_update_requests
WHERE user_id = $1 AND request_month = $2;

-- name: UpdateRequestStatus :exec
UPDATE face_profile_update_requests
SET status = $3, updated_at = CURRENT_TIMESTAMP
WHERE request_id = $1 AND company_id = $2;

-- name: ApproveRequest :exec
UPDATE face_profile_update_requests
SET 
    status = 1,
    approved_by = $3,
    approved_at = CURRENT_TIMESTAMP,
    update_token = $4,
    update_link_expires_at = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE request_id = $1 AND company_id = $2;

-- name: RejectRequest :exec
UPDATE face_profile_update_requests
SET 
    status = 2,
    approved_by = $3,
    rejection_reason = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE request_id = $1 AND company_id = $2;

-- name: CompleteRequest :exec
UPDATE face_profile_update_requests
SET status = 4, updated_at = CURRENT_TIMESTAMP
WHERE request_id = $1 AND company_id = $2;

-- name: MarkExpiredRequests :execrows
UPDATE face_profile_update_requests
SET status = 3, updated_at = CURRENT_TIMESTAMP
WHERE status = 1 
  AND update_link_expires_at IS NOT NULL 
  AND update_link_expires_at < CURRENT_TIMESTAMP;

-- name: HasPendingRequest :one
SELECT EXISTS(
    SELECT 1 
    FROM face_profile_update_requests 
    WHERE user_id = $1 AND company_id = $2 AND status = 0
) as has_pending;
