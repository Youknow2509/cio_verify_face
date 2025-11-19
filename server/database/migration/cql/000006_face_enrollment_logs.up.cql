CREATE TABLE IF NOT EXISTS face_enrollment_logs (
    company_id UUID,
    year_month TEXT,
    created_at TIMESTAMP,
    employee_id UUID,
    action_type TEXT,
    status TEXT,
    image_url TEXT,
    failure_reason TEXT,
    metadata MAP<TEXT, TEXT>,
    PRIMARY KEY ((company_id, year_month), created_at, employee_id)
) WITH CLUSTERING ORDER BY (created_at DESC);