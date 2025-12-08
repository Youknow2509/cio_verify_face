-- +goose Up
-- +goose StatementBegin

-- =================================================================
-- COMPANIES
-- =================================================================
INSERT INTO companies (name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices)
SELECT 'FPT Software', '17 Duy Tan, Cau Giay, Hanoi', '+84-24-7300-7300', 'contact@fpt.software', 'https://fpt-software.com', 1, 2, DATE '2025-01-01', DATE '2026-12-31', 10000, 200
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'FPT Software');

INSERT INTO companies (name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices)
SELECT 'VinGroup', '7 Bang Lang 1, Vinhomes Riverside, Long Bien, Hanoi', '+84-24-3974-9999', 'info@vingroup.net', 'https://www.vingroup.net', 1, 2, DATE '2025-03-15', DATE '2026-03-14', 50000, 1000
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'VinGroup');

INSERT INTO companies (name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices)
SELECT 'Techcombank', '191 Ba Trieu, Hai Ba Trung, Hanoi', '+84-1800-588-822', 'callcenter@techcombank.com.vn', 'https://techcombank.com', 1, 1, DATE '2025-02-01', DATE '2026-01-31', 12000, 500
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'Techcombank');

INSERT INTO companies (name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices)
SELECT 'Acme Tech', '123 Nguyen Trai, HCMC', '+84-28-0000-0001', 'contact@acmetech.example', 'https://acmetech.example', 1, 1, DATE '2025-01-01', DATE '2025-12-31', 500, 50
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'Acme Tech');

INSERT INTO companies (name, address, phone, email, website, status, subscription_plan, subscription_start_date, subscription_end_date, max_employees, max_devices)
SELECT 'Beta Retail', '456 Le Loi, HCMC', '+84-28-0000-0002', 'info@betaretail.example', 'https://betaretail.example', 1, 0, DATE '2025-02-01', DATE '2026-01-31', 200, 20
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'Beta Retail');

-- =================================================================
-- COMPANY SETTINGS
-- =================================================================
DO $$
DECLARE
    fpt_id UUID;
    vin_id UUID;
    tcb_id UUID;
    acme_id UUID;
    beta_id UUID;
BEGIN
    SELECT company_id INTO fpt_id FROM companies WHERE name = 'FPT Software';
    SELECT company_id INTO vin_id FROM companies WHERE name = 'VinGroup';
    SELECT company_id INTO tcb_id FROM companies WHERE name = 'Techcombank';
    SELECT company_id INTO acme_id FROM companies WHERE name = 'Acme Tech';
    SELECT company_id INTO beta_id FROM companies WHERE name = 'Beta Retail';

    -- FPT Settings
    INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
    VALUES
        (fpt_id, 'timezone', 'Asia/Ho_Chi_Minh', 0, 'Company timezone'),
        (fpt_id, 'workweek_start', '1', 1, '1=Monday'),
        (fpt_id, 'allow_remote_checkin', 'true', 2, 'Allow check-in from mobile app')
    ON CONFLICT (company_id, setting_key) DO NOTHING;

    -- VinGroup Settings
    INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
    VALUES
        (vin_id, 'timezone', 'Asia/Ho_Chi_Minh', 0, 'Company timezone'),
        (vin_id, 'overtime_policy', '{"rate": 1.5, "max_hours_per_month": 40}', 3, 'Overtime rules')
    ON CONFLICT (company_id, setting_key) DO NOTHING;

    -- Techcombank Settings
    INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
    VALUES
        (tcb_id, 'timezone', 'Asia/Ho_Chi_Minh', 0, 'Company timezone'),
        (tcb_id, 'security_level', 'high', 0, 'Security policy level')
    ON CONFLICT (company_id, setting_key) DO NOTHING;

    -- Acme Tech Settings
    INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
    VALUES
        (acme_id, 'timezone', 'Asia/Ho_Chi_Minh', 0, 'Company timezone'),
        (acme_id, 'workweek_start', '1', 1, '1=Monday')
    ON CONFLICT (company_id, setting_key) DO NOTHING;

    -- Beta Retail Settings
    INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description)
    VALUES
        (beta_id, 'timezone', 'Asia/Ho_Chi_Minh', 0, 'Company timezone')
    ON CONFLICT (company_id, setting_key) DO NOTHING;
END $$;


-- =================================================================
-- USERS & EMPLOYEES
-- =================================================================
-- This block creates users and then links them as employees to companies.
-- Passwords are the same as the email for simplicity.
DO $$
DECLARE
    fpt_id UUID;
    vin_id UUID;
    tcb_id UUID;
    acme_id UUID;
    beta_id UUID;
    user_id_var UUID;
BEGIN
    -- Get company IDs
    SELECT company_id INTO fpt_id FROM companies WHERE name = 'FPT Software';
    SELECT company_id INTO vin_id FROM companies WHERE name = 'VinGroup';
    SELECT company_id INTO tcb_id FROM companies WHERE name = 'Techcombank';
    SELECT company_id INTO acme_id FROM companies WHERE name = 'Acme Tech';
    SELECT company_id INTO beta_id FROM companies WHERE name = 'Beta Retail';

    -- FPT Software Users & Employees
    -- Admin
    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('admin.fpt@example.com', '0910000001', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'FPT Admin', 1, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, fpt_id, 'FPT0001', 'Management', 'Company Admin', '2023-01-10', 50000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;
    -- Employees
    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('employee1.fpt@example.com', '0910000002', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Nguyen Van A', 2, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, fpt_id, 'FPT0002', 'AI Lab', 'AI Engineer', '2024-03-01', 35000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('employee2.fpt@example.com', '0910000003', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Tran Thi B', 2, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, fpt_id, 'FPT0003', 'HR', 'HR Specialist', '2023-11-20', 28000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    -- VinGroup Users & Employees
    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('admin.vin@example.com', '0920000001', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'VinGroup Admin', 1, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, vin_id, 'VIN0001', 'Board of Directors', 'Director', '2022-05-15', 100000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    -- Techcombank Users & Employees
    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('admin.tcb@example.com', '0930000001', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Techcombank Admin', 1, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, tcb_id, 'TCB0001', 'IT', 'Branch Manager', '2023-08-01', 60000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    -- Acme Tech Users & Employees (from original file)
    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('admin.acme@example.com', '0900000001', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Acme Admin', 1, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, acme_id, 'ACM001', 'Engineering', 'Software Engineer', '2024-06-01', 20000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('alice.acme@example.com', '0900001001', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Alice Nguyen', 2, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, acme_id, 'ACM002', 'Engineering', 'Software Engineer', '2024-06-01', 20000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('bob.acme@example.com', '0900001002', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Bob Tran', 2, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, acme_id, 'ACM003', 'Operations', 'Ops Specialist', '2024-07-15', 15000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    -- Beta Retail Users & Employees (from original file)
    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('admin.beta@example.com', '0900000002', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Beta Admin', 1, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, beta_id, 'BTA001', 'Sales', 'Sales Associate', '2024-05-10', 12000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;

    INSERT INTO users (email, phone, salt, password_hash, full_name, role, status)
    VALUES ('charlie.beta@example.com', '0900002001', 'salt_test', 'd8d9e32681bae9ba7ae6ff47f216d9d978f22f52c99a34ef14baedcee67e82ee', 'Charlie Pham', 2, 0)
    ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name RETURNING user_id INTO user_id_var;
    IF user_id_var IS NOT NULL THEN
        INSERT INTO employees (employee_id, company_id, employee_code, department, position, hire_date, salary, status)
        VALUES (user_id_var, beta_id, 'BTA002', 'Sales', 'Sales Associate', '2024-05-10', 12000000, 0)
        ON CONFLICT (employee_id) DO NOTHING;
    END IF;
END $$;

-- =================================================================
-- DEVICES
-- =================================================================
DO $$
DECLARE
    fpt_id UUID;
    vin_id UUID;
    tcb_id UUID;
    acme_id UUID;
    beta_id UUID;
BEGIN
    SELECT company_id INTO fpt_id FROM companies WHERE name = 'FPT Software';
    SELECT company_id INTO vin_id FROM companies WHERE name = 'VinGroup';
    SELECT company_id INTO tcb_id FROM companies WHERE name = 'Techcombank';
    SELECT company_id INTO acme_id FROM companies WHERE name = 'Acme Tech';
    SELECT company_id INTO beta_id FROM companies WHERE name = 'Beta Retail';

    -- FPT Devices
    INSERT INTO devices (company_id, name, address, device_type, serial_number, mac_address, ip_address, status, token)
    VALUES
        (fpt_id, 'FPT Tower - Lobby A', '17 Duy Tan', 0, 'FPT-FT-001', '1A:2B:3C:4D:5E:01', '10.1.1.10'::inet, 1, 'device-token-fpt-1'),
        (fpt_id, 'FPT Tower - Canteen', '17 Duy Tan', 0, 'FPT-FT-002', '1A:2B:3C:4D:5E:02', '10.1.1.11'::inet, 1, 'device-token-fpt-2'),
        (fpt_id, 'F-Town 3 - Entrance', 'District 9, HCMC', 0, 'FPT-FT-003', '1A:2B:3C:4D:5E:03', '10.2.1.10'::inet, 0, 'device-token-fpt-3')
    ON CONFLICT (device_id) DO NOTHING;

    -- VinGroup Devices
    INSERT INTO devices (company_id, name, address, device_type, serial_number, mac_address, ip_address, status, token)
    VALUES
        (vin_id, 'Vincom Center - Main Gate', '72 Le Thanh Ton, HCMC', 0, 'VIN-FT-001', '2A:3B:4C:5D:6E:01', '192.168.10.100'::inet, 1, 'device-token-vin-1'),
        (vin_id, 'Vinhomes Central Park - Lobby', 'Binh Thanh, HCMC', 0, 'VIN-FT-002', '2A:3B:4C:5D:6E:02', '192.168.20.100'::inet, 1, 'device-token-vin-2')
    ON CONFLICT (device_id) DO NOTHING;

    -- Acme Tech Devices
    INSERT INTO devices (company_id, name, address, device_type, serial_number, mac_address, ip_address, status, token, settings)
    VALUES
        (acme_id, 'Face Terminal - Main Entrance', 'HQ Lobby', 0, 'ACME-FT-001', '00:11:22:33:44:55', '192.168.1.10'::inet, 1, 'device-token-acme-1', '{"camera":"v1"}'::jsonb),
        (acme_id, 'Webcam - Dev Room', 'Engineering Dept', 2, 'ACME-WC-001', '00:11:22:33:44:66', '192.168.1.20'::inet, 1, 'device-token-acme-2', '{}'::jsonb)
    ON CONFLICT (device_id) DO NOTHING;

    -- Beta Retail Devices
    INSERT INTO devices (company_id, name, address, device_type, serial_number, mac_address, ip_address, status, token)
    VALUES
        (beta_id, 'Mobile App - Sales Team', 'On-field staff', 1, 'BETA-MA-001', NULL, NULL, 1, 'device-token-beta-1')
    ON CONFLICT (device_id) DO NOTHING;
END $$;

-- =================================================================
-- WORK SHIFTS
-- =================================================================
DO $$
DECLARE
    fpt_id UUID;
    acme_id UUID;
    beta_id UUID;
BEGIN
    SELECT company_id INTO fpt_id FROM companies WHERE name = 'FPT Software';
    SELECT company_id INTO acme_id FROM companies WHERE name = 'Acme Tech';
    SELECT company_id INTO beta_id FROM companies WHERE name = 'Beta Retail';

    -- FPT Shifts
    INSERT INTO work_shifts (company_id, name, start_time, end_time, work_days, is_flexible)
    VALUES
        (fpt_id, 'Morning Shift', '08:00:00', '12:00:00', ARRAY[1,2,3,4,5], false),
        (fpt_id, 'Afternoon Shift', '13:30:00', '17:30:00', ARRAY[1,2,3,4,5], false),
        (fpt_id, 'Flexible Engineering', '09:00:00', '18:00:00', ARRAY[1,2,3,4,5], true)
    ON CONFLICT (shift_id) DO NOTHING;

    -- Acme Shifts
    INSERT INTO work_shifts (company_id, name, start_time, end_time, work_days)
    VALUES
        (acme_id, 'Standard Office Hours', '08:30:00', '17:30:00', ARRAY[1,2,3,4,5])
    ON CONFLICT (shift_id) DO NOTHING;

    -- Beta Shifts
    INSERT INTO work_shifts (company_id, name, start_time, end_time, work_days)
    VALUES
        (beta_id, 'Retail Shift A', '08:00:00', '16:00:00', ARRAY[1,2,3,4,5,6]),
        (beta_id, 'Retail Shift B', '14:00:00', '22:00:00', ARRAY[1,2,3,4,5,6,7])
    ON CONFLICT (shift_id) DO NOTHING;
END $$;

-- =================================================================
-- EMPLOYEE SHIFTS
-- =================================================================
DO $$
DECLARE
    emp_id UUID;
    shift_id_var UUID;
BEGIN
    -- Assign FPT Employee 1 to flexible shift
    SELECT user_id INTO emp_id FROM users WHERE email = 'employee1.fpt@example.com';
    SELECT shift_id INTO shift_id_var FROM work_shifts WHERE name = 'Flexible Engineering' AND company_id = (SELECT company_id FROM companies WHERE name = 'FPT Software');
    IF emp_id IS NOT NULL AND shift_id_var IS NOT NULL THEN
        INSERT INTO employee_shifts (employee_id, shift_id, effective_from)
        VALUES (emp_id, shift_id_var, '2024-01-01')
        ON CONFLICT (employee_id, shift_id, effective_from) DO NOTHING;
    END IF;

    -- Assign Acme employees to standard shift
    SELECT shift_id INTO shift_id_var FROM work_shifts WHERE name = 'Standard Office Hours' AND company_id = (SELECT company_id FROM companies WHERE name = 'Acme Tech');
    IF shift_id_var IS NOT NULL THEN
        -- Alice
        SELECT user_id INTO emp_id FROM users WHERE email = 'alice.acme@example.com';
        IF emp_id IS NOT NULL THEN
            INSERT INTO employee_shifts (employee_id, shift_id, effective_from)
            VALUES (emp_id, shift_id_var, '2024-06-01')
            ON CONFLICT (employee_id, shift_id, effective_from) DO NOTHING;
        END IF;
        -- Bob
        SELECT user_id INTO emp_id FROM users WHERE email = 'bob.acme@example.com';
        IF emp_id IS NOT NULL THEN
            INSERT INTO employee_shifts (employee_id, shift_id, effective_from)
            VALUES (emp_id, shift_id_var, '2024-07-15')
            ON CONFLICT (employee_id, shift_id, effective_from) DO NOTHING;
        END IF;
    END IF;

    -- Assign Beta employee to a retail shift
    SELECT user_id INTO emp_id FROM users WHERE email = 'charlie.beta@example.com';
    SELECT shift_id INTO shift_id_var FROM work_shifts WHERE name = 'Retail Shift A' AND company_id = (SELECT company_id FROM companies WHERE name = 'Beta Retail');
    IF emp_id IS NOT NULL AND shift_id_var IS NOT NULL THEN
        INSERT INTO employee_shifts (employee_id, shift_id, effective_from)
        VALUES (emp_id, shift_id_var, '2024-05-10')
        ON CONFLICT (employee_id, shift_id, effective_from) DO NOTHING;
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Truncate tables in reverse order of dependency to avoid foreign key constraints
TRUNCATE TABLE face_audit_logs, face_profiles, employee_shifts, work_shifts, devices, company_settings, system_settings, audit_logs, employees, user_sessions, users, companies RESTART IDENTITY CASCADE;
-- +goose StatementEnd
