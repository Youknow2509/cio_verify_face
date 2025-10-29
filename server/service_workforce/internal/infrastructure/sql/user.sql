-- name: CheckUserExistInCompany :one
SELECT 1
FROM employees
WHERE employee_id = $1 AND company_id = $2;

-- name: UserPermissionDevice :one
SELECT EXISTS (
    SELECT 1
    FROM devices d JOIN employees e ON d.company_id = e.company_id
WHERE e.employee_id = $1 AND d.device_id = $2
) AS exist;

