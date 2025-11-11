-- CREATE TABLE IF NOT EXISTS work_shifts (
-- shift_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
-- company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
-- name VARCHAR(255) NOT NULL,
-- description TEXT,
-- start_time TIME NOT NULL,
-- end_time TIME NOT NULL,
-- break_duration_minutes INTEGER DEFAULT 0,
-- grace_period_minutes INTEGER DEFAULT 15, -- Late arrival tolerance
-- early_departure_minutes INTEGER DEFAULT 15, -- Early leave tolerance
-- work_days INTEGER[] DEFAULT ARRAY[1,2,3,4,5], -- 1=Monday, 7=Sunday
-- is_flexible BOOLEAN DEFAULT FALSE,
-- overtime_after_minutes INTEGER DEFAULT 480, -- 8 hours = 480 minutes
-- is_active BOOLEAN DEFAULT TRUE,
-- created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
-- updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- )

-- name: GetListEmployeeInShift :many

-- List employees currently active in a specific shift (as of CURRENT_DATE)
SELECT 
        e.employee_id,
        u.full_name,
        e.employee_code,
        ws.name AS shift_name,
        TRUE AS current_shift
FROM employee_shifts es
INNER JOIN employees e ON e.employee_id = es.employee_id
INNER JOIN users u ON u.user_id = e.employee_id
INNER JOIN work_shifts ws ON ws.shift_id = es.shift_id
WHERE e.company_id = $1
    AND es.shift_id = $2
    AND es.is_active = TRUE
    AND es.effective_from <= CURRENT_DATE 
    AND (es.effective_to IS NULL OR es.effective_to >= CURRENT_DATE)
ORDER BY e.employee_code
LIMIT $3 OFFSET $4;

-- name: CountEmployeesInShiftCurrent :one
-- Count employees currently active in a specific shift (as of CURRENT_DATE)
SELECT COUNT(1) FROM (
        SELECT DISTINCT e.employee_id
        FROM employee_shifts es
        INNER JOIN employees e ON e.employee_id = es.employee_id
        WHERE e.company_id = $1
            AND es.shift_id = $2
            AND es.is_active = TRUE
            AND es.effective_from <= CURRENT_DATE
            AND (es.effective_to IS NULL OR es.effective_to >= CURRENT_DATE)
) t;

-- name: GetListEmployeeDonotInShift :many
-- List employees of a company who are NOT currently active in the given shift (as of CURRENT_DATE)
SELECT 
        e.employee_id,
        u.full_name,
        e.employee_code,
        COALESCE(ws.name, '') AS shift_name,
        FALSE AS current_shift
FROM employees e
INNER JOIN users u ON u.user_id = e.employee_id
LEFT JOIN LATERAL (
        SELECT es.shift_id
        FROM employee_shifts es
        WHERE es.employee_id = e.employee_id
            AND es.is_active = TRUE
            AND es.effective_from <= CURRENT_DATE
            AND (es.effective_to IS NULL OR es.effective_to >= CURRENT_DATE)
        ORDER BY es.effective_from DESC
        LIMIT 1
) ca ON TRUE
LEFT JOIN work_shifts ws ON ws.shift_id = ca.shift_id
WHERE e.company_id = $1
    AND NOT EXISTS (
        SELECT 1 FROM employee_shifts es2
        WHERE es2.employee_id = e.employee_id
            AND es2.shift_id = $2
            AND es2.is_active = TRUE
            AND es2.effective_from <= CURRENT_DATE
            AND (es2.effective_to IS NULL OR es2.effective_to >= CURRENT_DATE)
    )
ORDER BY e.employee_code
LIMIT $3 OFFSET $4;

-- name: CountEmployeesDonotInShiftCurrent :one
-- Count employees of a company who are NOT currently active in the given shift (as of CURRENT_DATE)
SELECT COUNT(1) FROM employees e
WHERE e.company_id = $1
    AND NOT EXISTS (
        SELECT 1 FROM employee_shifts es
        WHERE es.employee_id = e.employee_id
            AND es.shift_id = $2
            AND es.is_active = TRUE
            AND es.effective_from <= CURRENT_DATE
            AND (es.effective_to IS NULL OR es.effective_to >= CURRENT_DATE)
    );

-- name: IsUserManagetShift :one
SELECT shift_id 
FROM work_shifts
WHERE shift_id = $1 AND company_id = $2
LIMIT 1;

-- name: CreateShift :one
INSERT INTO work_shifts (
  company_id, name, description, start_time, end_time,
  break_duration_minutes, grace_period_minutes, early_departure_minutes,
  work_days, is_active
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8,
    $9, TRUE
)
RETURNING shift_id;

-- name: ListShifts :many
SELECT *
FROM work_shifts
WHERE company_id = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: GetShiftByID :one
SELECT *
FROM work_shifts
WHERE shift_id = $1;

-- name: UpdateTimeShift :exec
UPDATE work_shifts
SET start_time = $2,
    end_time = $3,
    break_duration_minutes = $4,
    grace_period_minutes = $5,
    early_departure_minutes = $6,
    work_days = $7,
    updated_at = now()
WHERE shift_id = $1;

-- name: DeleteShift :exec
DELETE FROM work_shifts
WHERE shift_id = $1;

-- name: GetShiftsIdForCompany :many
SELECT shift_id
FROM work_shifts
WHERE company_id = $1
LIMIT $2 OFFSET $3;

-- name: DisableShiftWithId :exec
UPDATE work_shifts
SET is_active = FALSE,
    updated_at = now()
WHERE shift_id = $1 and company_id = $2;

-- name: EnableShiftWithId :exec
UPDATE work_shifts
SET is_active = TRUE,
    updated_at = now()
WHERE shift_id = $1 and company_id = $2;

-- name: GetInfoEmployeeInShift :many
SELECT 
    e.employee_id,
    e.employee_code,
    u.full_name,
    CASE 
        WHEN es.shift_id = $2 THEN TRUE 
        ELSE FALSE 
    END AS current_shift,
    ws.name AS shift_active
FROM employees e
INNER JOIN users u ON e.employee_id = u.user_id
LEFT JOIN employee_shifts es ON e.employee_id = es.employee_id 
    AND es.is_active = TRUE
    AND es.effective_from <= CURRENT_DATE 
    AND (es.effective_to IS NULL OR es.effective_to >= CURRENT_DATE)
LEFT JOIN work_shifts ws ON es.shift_id = ws.shift_id
WHERE e.company_id = $1
ORDER BY e.employee_code;

