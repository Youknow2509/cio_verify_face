-- Table: employee_shifts
-- employee_shift_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
-- employee_id UUID NOT NULL REFERENCES employees(employee_id) ON DELETE CASCADE,
-- shift_id UUID NOT NULL REFERENCES work_shifts(shift_id) ON DELETE CASCADE,
-- effective_from DATE NOT NULL,
-- effective_to DATE,
-- is_active BOOLEAN DEFAULT TRUE,
-- created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
-- UNIQUE(employee_id, shift_id, effective_from, effective_to)

-- name: EditEffectiveShiftForEmployee :exec
UPDATE employee_shifts
SET effective_from = $2,
    effective_to = $3
WHERE employee_shift_id = $1;

-- name: DeleteEmployeeShift :exec
DELETE FROM employee_shifts
WHERE employee_shift_id = $1;

-- name: DisableEmployeeShift :exec
UPDATE employee_shifts
SET is_active = false
WHERE employee_shift_id = $1;

-- name: EnableEmployeeShift :exec
UPDATE employee_shifts
SET is_active = true
WHERE employee_shift_id = $1;

-- name: AddShiftForEmployee :exec
INSERT INTO employee_shifts (
    employee_id,
    shift_id,
    effective_from,
    effective_to
) VALUES ($1, $2, $3, $4)
RETURNING employee_shift_id;

-- name: CheckUserExistShift :one
SELECT employee_id,
    effective_from,
    effective_to
FROM employee_shifts
WHERE employee_id = $1 AND
    effective_from <= $2 AND
    effective_to >= $3 AND
    is_active = true
LIMIT $4 OFFSET $5;
