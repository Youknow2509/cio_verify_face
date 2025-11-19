CREATE TABLE IF NOT EXISTS audit_logs (
    company_id UUID,
    year_month TEXT,
    created_at TIMESTAMP,
    actor_id UUID,
    action_category TEXT,
    action_name TEXT,
    resource_type TEXT,
    resource_id TEXT,
    details MAP<TEXT, TEXT>,
    ip_address TEXT,
    user_agent TEXT,
    status TEXT,
    PRIMARY KEY ((company_id, year_month), created_at, actor_id)
) WITH CLUSTERING ORDER BY (created_at DESC);