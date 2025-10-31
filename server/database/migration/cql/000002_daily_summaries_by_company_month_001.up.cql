-- Daily summaries partitioned by company and month
CREATE TABLE IF NOT EXISTS daily_summaries_by_company_month (
    company_id UUID,
    summary_month TEXT, -- 'YYYY-MM' format
    work_date DATE,
    employee_id UUID,
    shift_id UUID,
    scheduled_in TIME,
    scheduled_out TIME,
    actual_check_in TIMESTAMP,
    actual_check_out TIMESTAMP,
    total_work_minutes INT,
    break_minutes INT,
    overtime_minutes INT,
    late_minutes INT,
    early_leave_minutes INT,
    attendance_status INT, -- 0: PRESENT, 1: LATE, 2: EARLY_LEAVE, 3: ABSENT
    attendance_percentage DECIMAL,
    notes TEXT,
    approved_by UUID,
    approved_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY ((company_id, summary_month), work_date, employee_id)
) WITH CLUSTERING ORDER BY (work_date DESC)
AND compression = {'class': 'LZ4Compressor'}
AND compaction = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_unit': 'DAYS',
    'compaction_window_size': 7
}
AND default_time_to_live = 31536000; -- 1 year TTL

-- Index for employee-based queries
CREATE INDEX IF NOT EXISTS daily_summary_employee_idx 
ON daily_summaries_by_company_month (employee_id);