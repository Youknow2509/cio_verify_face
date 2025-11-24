CREATE TABLE IF NOT EXISTS daily_summaries_by_user (
    company_id UUID,
    employee_id UUID,
    summary_month TEXT,    
    work_date DATE,
    shift_id UUID,
    actual_check_in TIMESTAMP,
    actual_check_out TIMESTAMP,
    attendance_status INT,
    late_minutes INT,
    early_leave_minutes INT,
    total_work_minutes INT,
    notes TEXT,
    updated_at TIMESTAMP,
    PRIMARY KEY ((company_id, employee_id, summary_month), work_date)
) WITH CLUSTERING ORDER BY (work_date DESC);