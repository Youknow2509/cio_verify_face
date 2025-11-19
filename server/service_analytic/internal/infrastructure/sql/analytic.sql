-- name: GetDailyAttendanceSummaryByDate :many
SELECT 
    summary_id,
    employee_id,
    shift_id,
    work_date,
    scheduled_in,
    scheduled_out,
    actual_check_in,
    actual_check_out,
    total_work_minutes,
    break_minutes,
    overtime_minutes,
    late_minutes,
    early_leave_minutes,
    status,
    attendance_percentage,
    notes,
    approved_by,
    approved_at,
    created_at,
    updated_at
FROM daily_attendance_summary
WHERE work_date = $1
AND ($2::uuid IS NULL OR employee_id IN (
    SELECT employee_id FROM employees WHERE company_id = $2
))
ORDER BY employee_id;

-- name: GetAttendanceSummaryByDateRange :many
SELECT 
    summary_id,
    employee_id,
    shift_id,
    work_date,
    scheduled_in,
    scheduled_out,
    actual_check_in,
    actual_check_out,
    total_work_minutes,
    break_minutes,
    overtime_minutes,
    late_minutes,
    early_leave_minutes,
    status,
    attendance_percentage,
    notes,
    approved_by,
    approved_at,
    created_at,
    updated_at
FROM daily_attendance_summary
WHERE work_date >= $1 AND work_date <= $2
AND ($3::uuid IS NULL OR employee_id IN (
    SELECT employee_id FROM employees WHERE company_id = $3
))
ORDER BY work_date, employee_id;

-- name: GetEmployeeByID :one
SELECT 
    employee_id,
    company_id,
    employee_code,
    department,
    position,
    hire_date,
    salary,
    status,
    created_at,
    updated_at
FROM employees
WHERE employee_id = $1
LIMIT 1;

-- name: GetEmployeesByCompany :many
SELECT 
    employee_id,
    company_id,
    employee_code,
    department,
    position,
    hire_date,
    salary,
    status,
    created_at,
    updated_at
FROM employees
WHERE company_id = $1
ORDER BY employee_code;

-- name: GetTotalEmployeesCount :one
SELECT COUNT(*) as count
FROM employees
WHERE ($1::uuid IS NULL OR company_id = $1)
AND status = 0;

-- name: GetUserByID :one
SELECT 
    user_id,
    full_name,
    email,
    role
FROM users
WHERE user_id = $1
LIMIT 1;

-- name: GetWorkShiftByID :one
SELECT 
    shift_id,
    company_id,
    name,
    start_time::text as start_time,
    end_time::text as end_time
FROM work_shifts
WHERE shift_id = $1
LIMIT 1;

-- name: GetWorkShiftsByCompany :many
SELECT 
    shift_id,
    company_id,
    name,
    start_time::text as start_time,
    end_time::text as end_time
FROM work_shifts
WHERE company_id = $1
ORDER BY start_time;

-- name: GetCompanyByID :one
SELECT 
    company_id,
    name,
    address
FROM companies
WHERE company_id = $1
LIMIT 1;

-- name: GetAttendanceRecordsByDateRange :many
SELECT 
    record_id,
    employee_id,
    device_id,
    timestamp,
    record_type,
    verification_method,
    verification_score,
    face_image_url,
    sync_status,
    created_at
FROM attendance_records
WHERE employee_id = $1
AND timestamp >= $2
AND timestamp <= $3
ORDER BY timestamp;
