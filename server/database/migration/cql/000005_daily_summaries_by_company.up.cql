CREATE TABLE IF NOT EXISTS daily_summaries_by_company (
    company_id UUID,
    summary_month TEXT,    
    work_date DATE,
    employee_id UUID,
    shift_id UUID,
    actual_check_in TIMESTAMP,
    actual_check_out TIMESTAMP,
    attendance_status INT,   
    late_minutes INT,
    early_leave_minutes INT,
    total_work_minutes INT,
    notes TEXT,
    updated_at TIMESTAMP,
    PRIMARY KEY ((company_id, summary_month), work_date, employee_id)
) WITH CLUSTERING ORDER BY (work_date DESC, employee_id ASC);