-- =================================================================
-- USER QUERIES
-- =================================================================

-- name: GetUserByID :one
SELECT 
    user_id,
    email,
    phone,
    full_name,
    avatar_url,
    role,
    status,
    is_locked,
    created_at,
    updated_at
FROM users
WHERE user_id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT 
    user_id,
    email,
    phone,
    full_name,
    avatar_url,
    role,
    status,
    is_locked,
    created_at,
    updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: CheckUserBelongsToCompany :one
SELECT EXISTS(
    SELECT 1 
    FROM employees 
    WHERE employee_id = $1 AND company_id = $2
) as belongs;

-- name: GetEmployeeInfo :one
SELECT 
    e.employee_id,
    e.company_id,
    e.employee_code,
    e.department,
    e.position,
    e.hire_date,
    e.status,
    u.email,
    u.phone,
    u.full_name,
    u.avatar_url
FROM employees e
JOIN users u ON e.employee_id = u.user_id
WHERE e.employee_id = $1
LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET 
    salt = $2,
    password_hash = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: IsCompanyAdmin :one
SELECT EXISTS(
    SELECT 1 
    FROM employees 
    WHERE employee_id = $1 
      AND company_id = $2 
      AND employee_id IN (
          SELECT user_id 
          FROM users 
          WHERE role = 1
      )
) as is_admin;
