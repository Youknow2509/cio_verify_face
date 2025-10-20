-- name: CheckUserExistInCompany :one
SELECT 1
FROM employees
WHERE employee_id = $1 AND company_id = $2;
