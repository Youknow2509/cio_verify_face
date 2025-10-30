-- +goose Up
-- +goose StatementBegin
-- Companies
INSERT INTO companies (name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices)
SELECT 'Acme Tech', '123 Nguyen Trai, HCMC', '+84-28-0000-0001', 'contact@acmetech.example', 'https://acmetech.example', 1, 1,
			 DATE '2025-01-01', DATE '2025-12-31', 500, 50
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'Acme Tech');

INSERT INTO companies (name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices)
SELECT 'Beta Retail', '456 Le Loi, HCMC', '+84-28-0000-0002', 'info@betaretail.example', 'https://betaretail.example', 1, 0,
			 DATE '2025-02-01', DATE '2026-01-31', 200, 20
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'Beta Retail');

-- Company settings
INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
SELECT c.company_id, 'timezone', 'Asia/Ho_Chi_Minh', 0, 'Company timezone'
FROM companies c
WHERE c.name = 'Acme Tech'
	AND NOT EXISTS (
		SELECT 1 FROM company_settings cs WHERE cs.company_id = c.company_id AND cs.setting_key = 'timezone'
	);

INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
SELECT c.company_id, 'workweek_start', '1', 1, '1=Monday'
FROM companies c
WHERE c.name = 'Acme Tech'
	AND NOT EXISTS (
		SELECT 1 FROM company_settings cs WHERE cs.company_id = c.company_id AND cs.setting_key = 'workweek_start'
	);

INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
SELECT c.company_id, 'timezone', 'Asia/Ho_Chi_Minh', 0, 'Company timezone'
FROM companies c
WHERE c.name = 'Beta Retail'
	AND NOT EXISTS (
		SELECT 1 FROM company_settings cs WHERE cs.company_id = c.company_id AND cs.setting_key = 'timezone'
	);

-- Users: 2 company admins + 3 employees
-- Password hash strings here are placeholders for tests only 
-- Password same email
INSERT INTO users (email, phone, salt, password_hash, full_name, avatar_url, role, status, is_locked)
SELECT 'admin.acme@example.com', '0900000001', 'salt_test', '38d6553f30f9c36131bb86d18055d013c0ca8dfc785aa93217b3eaadf5543f9c', 'Acme Admin', NULL, 1, 0, FALSE
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin.acme@example.com');

INSERT INTO users (email, phone, salt, password_hash, full_name, avatar_url, role, status, is_locked)
SELECT 'alice.acme@example.com', '0900001001', 'salt_test', '4c221a595c3ad4dac48865406408ecaea0cd48f23236cbf5bfbf545f8e642d44', 'Alice Nguyen', NULL, 2, 0, FALSE
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'alice.acme@example.com');

INSERT INTO users (email, phone, salt, password_hash, full_name, avatar_url, role, status, is_locked)
SELECT 'bob.acme@example.com', '0900001002', 'salt_test', 'c2c7d6d7eb4b098f9d627097ece3c013cb58b5e8998dcc703795f2e975a7feb2', 'Bob Tran', NULL, 2, 0, FALSE
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'bob.acme@example.com');

INSERT INTO users (email, phone, salt, password_hash, full_name, avatar_url, role, status, is_locked)
SELECT 'admin.beta@example.com', '0900000002', 'salt_test', '6240f82c69bff33493d351e7d7c028a603f803f10292b5b5e6d97e1fc47eac80', 'Beta Admin', NULL, 1, 0, FALSE
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin.beta@example.com');

INSERT INTO users (email, phone, salt, password_hash, full_name, avatar_url, role, status, is_locked)
SELECT 'charlie.beta@example.com', '0900002001', 'salt_test', '1a7ba7c24e0314e9b9da22c03ac973bee87fb28ca26caf9843046efb270984c7', 'Charlie Pham', NULL, 2, 0, FALSE
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'charlie.beta@example.com');

-- Employees
INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
SELECT u.user_id,
             (SELECT company_id FROM companies WHERE name = 'Acme Tech'),
             'ACM001', 'Engineering', 'Software Engineer', DATE '2024-06-01', 20000000.00, 0
FROM users u WHERE u.email = 'admin.acme@example.com'
ON CONFLICT (employee_id) DO NOTHING;

INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
SELECT u.user_id,
			 (SELECT company_id FROM companies WHERE name = 'Acme Tech'),
			 'ACM001', 'Engineering', 'Software Engineer', DATE '2024-06-01', 20000000.00, 0
FROM users u WHERE u.email = 'alice.acme@example.com'
ON CONFLICT (employee_id) DO NOTHING;

INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
SELECT u.user_id,
			 (SELECT company_id FROM companies WHERE name = 'Acme Tech'),
			 'ACM002', 'Operations', 'Ops Specialist', DATE '2024-07-15', 15000000.00, 0
FROM users u WHERE u.email = 'bob.acme@example.com'
ON CONFLICT (employee_id) DO NOTHING;

INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
SELECT u.user_id,
             (SELECT company_id FROM companies WHERE name = 'Beta Retail'),
             'BTA001', 'Sales', 'Sales Associate', DATE '2024-05-10', 12000000.00, 0
FROM users u WHERE u.email = 'admin.beta@example.com'
ON CONFLICT (employee_id) DO NOTHING;

INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
SELECT u.user_id,
			 (SELECT company_id FROM companies WHERE name = 'Beta Retail'),
			 'BTA001', 'Sales', 'Sales Associate', DATE '2024-05-10', 12000000.00, 0
FROM users u WHERE u.email = 'charlie.beta@example.com'
ON CONFLICT (employee_id) DO NOTHING;

-- Devices per company
INSERT INTO devices (company_id, name, address, device_type, serial_number, mac_address, ip_address, firmware_version, status, token, settings)
SELECT c.company_id, 'Face Terminal - Main Entrance', 'HQ Lobby', 0, 'ACME-FT-001', '00:11:22:33:44:55', '192.168.1.10'::inet, '1.0.0', 1, 'device-token-acme-1', '{"camera":"v1"}'::jsonb
FROM companies c
WHERE c.name = 'Acme Tech'
	AND NOT EXISTS (
		SELECT 1 FROM devices d WHERE d.company_id = c.company_id AND d.serial_number = 'ACME-FT-001'
	);

INSERT INTO devices (company_id, name, address, device_type, serial_number, mac_address, ip_address, firmware_version, status, token, settings)
SELECT c.company_id, 'Face Terminal - Store 1', 'Store 1 Entrance', 0, 'BETA-FT-001', 'AA:BB:CC:DD:EE:FF', '10.0.0.10'::inet, '1.0.0', 1, 'device-token-beta-1', '{"camera":"v1"}'::jsonb
FROM companies c
WHERE c.name = 'Beta Retail'
	AND NOT EXISTS (
		SELECT 1 FROM devices d WHERE d.company_id = c.company_id AND d.serial_number = 'BETA-FT-001'
	);

-- Work shifts
INSERT INTO work_shifts (company_id, name, description, start_time, end_time, break_duration_minutes, grace_period_minutes, early_departure_minutes, work_days, is_flexible, overtime_after_minutes, is_active)
SELECT c.company_id, 'Day Shift', '9 AM - 6 PM', TIME '09:00', TIME '18:00', 60, 15, 15, ARRAY[1,2,3,4,5], FALSE, 480, TRUE
FROM companies c
WHERE c.name = 'Acme Tech'
	AND NOT EXISTS (
		SELECT 1 FROM work_shifts ws WHERE ws.company_id = c.company_id AND ws.name = 'Day Shift'
	);

INSERT INTO work_shifts (company_id, name, description, start_time, end_time, break_duration_minutes, grace_period_minutes, early_departure_minutes, work_days, is_flexible, overtime_after_minutes, is_active)
SELECT c.company_id, 'Night Shift', '10 PM - 6 AM', TIME '22:00', TIME '06:00', 30, 10, 10, ARRAY[1,2,3,4,5,6], FALSE, 480, TRUE
FROM companies c
WHERE c.name = 'Acme Tech'
	AND NOT EXISTS (
		SELECT 1 FROM work_shifts ws WHERE ws.company_id = c.company_id AND ws.name = 'Night Shift'
	);

INSERT INTO work_shifts (company_id, name, description, start_time, end_time, break_duration_minutes, grace_period_minutes, early_departure_minutes, work_days, is_flexible, overtime_after_minutes, is_active)
SELECT c.company_id, 'Retail Shift', '8 AM - 5 PM', TIME '08:00', TIME '17:00', 60, 10, 10, ARRAY[1,2,3,4,5,6], FALSE, 480, TRUE
FROM companies c
WHERE c.name = 'Beta Retail'
	AND NOT EXISTS (
		SELECT 1 FROM work_shifts ws WHERE ws.company_id = c.company_id AND ws.name = 'Retail Shift'
	);

-- Assign employee shifts (effective from Oct 1, 2025)
INSERT INTO employee_shifts (employee_id, shift_id, effective_from, effective_to, is_active)
SELECT e.employee_id,
			 (SELECT shift_id FROM work_shifts ws JOIN companies c ON ws.company_id=c.company_id WHERE c.name='Acme Tech' AND ws.name='Day Shift'),
			 DATE '2025-10-01', NULL, TRUE
FROM employees e
JOIN users u ON u.user_id = e.employee_id
JOIN companies c ON c.company_id = e.company_id
WHERE c.name='Acme Tech' AND u.email IN ('alice.acme@example.com','bob.acme@example.com')
ON CONFLICT DO NOTHING;

INSERT INTO employee_shifts (employee_id, shift_id, effective_from, effective_to, is_active)
SELECT e.employee_id,
			 (SELECT shift_id FROM work_shifts ws JOIN companies c ON ws.company_id=c.company_id WHERE c.name='Beta Retail' AND ws.name='Retail Shift'),
			 DATE '2025-10-01', NULL, TRUE
FROM employees e
JOIN users u ON u.user_id = e.employee_id
JOIN companies c ON c.company_id = e.company_id
WHERE c.name='Beta Retail' AND u.email IN ('charlie.beta@example.com')
ON CONFLICT DO NOTHING;

-- Attendance records for 2025-10-28 and 2025-10-29
-- Alice (Acme) - a bit late on 28th, on time on 29th
INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, face_image_url, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Acme Tech' AND d.serial_number='ACME-FT-001'),
			 TIMESTAMPTZ '2025-10-28 09:03:00+07', 0, 'FACE', 0.972, NULL, '{"temp":36.7}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='alice.acme@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, face_image_url, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Acme Tech' AND d.serial_number='ACME-FT-001'),
			 TIMESTAMPTZ '2025-10-28 18:10:00+07', 1, 'FACE', 0.981, NULL, '{"temp":36.6}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='alice.acme@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, face_image_url, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Acme Tech' AND d.serial_number='ACME-FT-001'),
			TIMESTAMPTZ '2025-10-29 09:00:00+07', 0, 'FACE', 0.965, NULL, '{}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='alice.acme@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, face_image_url, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Acme Tech' AND d.serial_number='ACME-FT-001'),
			TIMESTAMPTZ '2025-10-29 18:05:00+07', 1, 'FACE', 0.979, NULL, '{}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='alice.acme@example.com'
ON CONFLICT DO NOTHING;

-- Bob (Acme) - on time both days
INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Acme Tech' AND d.serial_number='ACME-FT-001'),
			TIMESTAMPTZ '2025-10-28 09:00:00+07', 0, 'FACE', 0.957, '{}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='bob.acme@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Acme Tech' AND d.serial_number='ACME-FT-001'),
			TIMESTAMPTZ '2025-10-28 18:00:00+07', 1, 'FACE', 0.962, '{}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='bob.acme@example.com'
ON CONFLICT DO NOTHING;

-- Charlie (Beta) - one day only
INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Beta Retail' AND d.serial_number='BETA-FT-001'),
			TIMESTAMPTZ '2025-10-28 08:00:00+07', 0, 'FACE', 0.953, '{}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='charlie.beta@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO attendance_records (employee_id, device_id, timestamp, record_type, verification_method, verification_score, metadata, sync_status)
SELECT e.employee_id,
			 (SELECT device_id FROM devices d JOIN companies c ON d.company_id=c.company_id WHERE c.name='Beta Retail' AND d.serial_number='BETA-FT-001'),
			TIMESTAMPTZ '2025-10-28 17:00:00+07', 1, 'FACE', 0.959, '{}'::jsonb, 0
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='charlie.beta@example.com'
ON CONFLICT DO NOTHING;

-- Daily attendance summary (derived-style examples)
-- Alice summaries
INSERT INTO daily_attendance_summary (employee_id, shift_id, work_date, scheduled_in, scheduled_out, actual_check_in, actual_check_out, total_work_minutes, break_minutes, overtime_minutes, late_minutes, early_leave_minutes, status, attendance_percentage, notes)
SELECT e.employee_id,
			 (SELECT shift_id FROM work_shifts ws JOIN companies c ON ws.company_id=c.company_id WHERE c.name='Acme Tech' AND ws.name='Day Shift'),
			 DATE '2025-10-28', TIME '09:00', TIME '18:00', TIMESTAMPTZ '2025-10-28 09:03:00+07', TIMESTAMPTZ '2025-10-28 18:10:00+07', 610, 60, 10, 3, 0, 1, 100.00, 'Test summary - Alice 28'
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='alice.acme@example.com'
ON CONFLICT (employee_id, work_date) DO NOTHING;

INSERT INTO daily_attendance_summary (employee_id, shift_id, work_date, scheduled_in, scheduled_out, actual_check_in, actual_check_out, total_work_minutes, break_minutes, overtime_minutes, late_minutes, early_leave_minutes, status, attendance_percentage, notes)
SELECT e.employee_id,
			 (SELECT shift_id FROM work_shifts ws JOIN companies c ON ws.company_id=c.company_id WHERE c.name='Acme Tech' AND ws.name='Day Shift'),
			 DATE '2025-10-29', TIME '09:00', TIME '18:00', TIMESTAMPTZ '2025-10-29 09:00:00+07', TIMESTAMPTZ '2025-10-29 18:05:00+07', 605, 60, 5, 0, 0, 0, 100.00, 'Test summary - Alice 29'
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='alice.acme@example.com'
ON CONFLICT (employee_id, work_date) DO NOTHING;

-- Bob summary
INSERT INTO daily_attendance_summary (employee_id, shift_id, work_date, scheduled_in, scheduled_out, actual_check_in, actual_check_out, total_work_minutes, break_minutes, overtime_minutes, late_minutes, early_leave_minutes, status, attendance_percentage, notes)
SELECT e.employee_id,
			 (SELECT shift_id FROM work_shifts ws JOIN companies c ON ws.company_id=c.company_id WHERE c.name='Acme Tech' AND ws.name='Day Shift'),
			 DATE '2025-10-28', TIME '09:00', TIME '18:00', TIMESTAMPTZ '2025-10-28 09:00:00+07', TIMESTAMPTZ '2025-10-28 18:00:00+07', 600, 60, 0, 0, 0, 0, 100.00, 'Test summary - Bob 28'
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='bob.acme@example.com'
ON CONFLICT (employee_id, work_date) DO NOTHING;

-- Charlie summary
INSERT INTO daily_attendance_summary (employee_id, shift_id, work_date, scheduled_in, scheduled_out, actual_check_in, actual_check_out, total_work_minutes, break_minutes, overtime_minutes, late_minutes, early_leave_minutes, status, attendance_percentage, notes)
SELECT e.employee_id,
			 (SELECT shift_id FROM work_shifts ws JOIN companies c ON ws.company_id=c.company_id WHERE c.name='Beta Retail' AND ws.name='Retail Shift'),
			 DATE '2025-10-28', TIME '08:00', TIME '17:00', TIMESTAMPTZ '2025-10-28 08:00:00+07', TIMESTAMPTZ '2025-10-28 17:00:00+07', 540, 60, 0, 0, 0, 0, 100.00, 'Test summary - Charlie 28'
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='charlie.beta@example.com'
ON CONFLICT (employee_id, work_date) DO NOTHING;

-- Attendance exception (late excuse for Alice on 28th)
INSERT INTO attendance_exceptions (summary_id, exception_type, reason, requested_by, approved_by, status, adjustment_minutes)
SELECT s.summary_id, 0, 'Test: Late due to traffic',
			 (SELECT user_id FROM users WHERE email='alice.acme@example.com'),
			 (SELECT user_id FROM users WHERE email='admin.acme@example.com'),
			 1, 3
FROM daily_attendance_summary s
JOIN employees e ON e.employee_id = s.employee_id
JOIN users u ON u.user_id = e.employee_id
WHERE u.email='alice.acme@example.com' AND s.work_date = DATE '2025-10-28'
ON CONFLICT DO NOTHING;

-- User session for Acme admin
INSERT INTO user_sessions (user_id, refresh_token, ip_address, user_agent, expires_at, is_active)
SELECT (SELECT user_id FROM users WHERE email='admin.acme@example.com'), 'test-refresh-token-acme-admin', '203.0.113.10', 'Mozilla/5.0 (Test) AppleWebKit', NOW() + INTERVAL '7 days', TRUE
WHERE NOT EXISTS (
	SELECT 1 FROM user_sessions WHERE user_id = (SELECT user_id FROM users WHERE email='admin.acme@example.com') AND is_active = TRUE
);

-- Audit log sample
INSERT INTO audit_logs (user_id, action, resource_type, resource_id, old_values, new_values, ip_address, user_agent)
SELECT (SELECT user_id FROM users WHERE email='admin.acme@example.com'), 'TEST_CREATE_EMPLOYEE', 'EMPLOYEE', e.employee_id, NULL, '{"created":true}'::jsonb, '203.0.113.10', 'Mozilla/5.0 (Test)'
FROM employees e JOIN users u ON u.user_id=e.employee_id WHERE u.email='alice.acme@example.com'
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove audit logs created by this seed
DELETE FROM audit_logs WHERE action = 'TEST_CREATE_EMPLOYEE';

-- Remove sessions for our test users
DELETE FROM user_sessions WHERE user_id IN (
	SELECT user_id FROM users WHERE email IN (
		'admin.acme@example.com','alice.acme@example.com','bob.acme@example.com','admin.beta@example.com','charlie.beta@example.com'
	)
);

-- Remove attendance exceptions and summaries for specific users/dates
DELETE FROM attendance_exceptions WHERE summary_id IN (
	SELECT summary_id FROM daily_attendance_summary s
	JOIN users u ON u.user_id = s.employee_id
	WHERE u.email IN ('alice.acme@example.com','bob.acme@example.com','charlie.beta@example.com')
		AND s.work_date IN (DATE '2025-10-28', DATE '2025-10-29')
);

DELETE FROM daily_attendance_summary
USING users u
WHERE u.user_id = daily_attendance_summary.employee_id
	AND u.email IN ('alice.acme@example.com','bob.acme@example.com','charlie.beta@example.com')
	AND work_date IN (DATE '2025-10-28', DATE '2025-10-29');

-- Remove attendance records for the test users
DELETE FROM attendance_records ar
USING users u
WHERE u.user_id = ar.employee_id
	AND u.email IN ('alice.acme@example.com','bob.acme@example.com','charlie.beta@example.com');

-- Remove employee shift assignments
DELETE FROM employee_shifts es
USING users u
WHERE u.user_id = es.employee_id
	AND u.email IN ('alice.acme@example.com','bob.acme@example.com','charlie.beta@example.com');

-- Remove devices tied to test companies
DELETE FROM devices d USING companies c
WHERE d.company_id = c.company_id AND c.name IN ('Acme Tech','Beta Retail');

-- Remove company settings for test companies
DELETE FROM company_settings cs USING companies c
WHERE cs.company_id = c.company_id AND c.name IN ('Acme Tech','Beta Retail');

-- Remove work shifts for test companies
DELETE FROM work_shifts ws USING companies c
WHERE ws.company_id = c.company_id AND c.name IN ('Acme Tech','Beta Retail');

-- Remove employees by user email (cascade from deleting users also works, this is defensive)
DELETE FROM employees e USING users u
WHERE e.employee_id = u.user_id AND u.email IN (
	'alice.acme@example.com','bob.acme@example.com','charlie.beta@example.com'
);

-- Remove users we created
DELETE FROM users WHERE email IN (
	'admin.acme@example.com','alice.acme@example.com','bob.acme@example.com','admin.beta@example.com','charlie.beta@example.com'
);

-- Remove companies we created
DELETE FROM companies WHERE name IN ('Acme Tech','Beta Retail');

-- Remove system settings we created if not desired
DELETE FROM system_settings WHERE setting_key IN ('password_min_length','allow_self_signup','default_timezone');
-- +goose StatementEnd
