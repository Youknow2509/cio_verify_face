-- name: GetIdCompanyByUserId :one
SELECT company_id
FROM employees
WHERE employee_id = $1
LIMIT 1;

-- name: CheckUserExistInCompany :one
SELECT 1
FROM employees
WHERE employee_id = $1 AND company_id = $2;

-- name: UserPermissionDevice :one
SELECT e.employee_id
FROM devices d JOIN employees e ON d.company_id = e.company_id
WHERE e.employee_id = $1 AND d.device_id = $2
LIMIT 1;

