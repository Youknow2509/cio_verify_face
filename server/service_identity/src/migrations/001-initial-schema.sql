-- Create companies table
CREATE TABLE IF NOT EXISTS companies (
    company_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    email VARCHAR(100),
    phone VARCHAR(20),
    status VARCHAR(50) DEFAULT 'ACTIVE',
    plan VARCHAR(100) DEFAULT 'FREE',
    max_employees INTEGER DEFAULT 0,
    settings JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create companies_secret table
CREATE TABLE IF NOT EXISTS companies_secret (
    company_secret_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    salt VARCHAR(255) NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    UNIQUE(company_id)
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    salt VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role INTEGER DEFAULT 2, -- 0: SYSTEM_ADMIN, 1: COMPANY_ADMIN, 2: EMPLOYEE
    is_active BOOLEAN DEFAULT TRUE,
    face_registered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create employees table
CREATE TABLE IF NOT EXISTS employees (
    employee_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    user_id UUID UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    employee_code VARCHAR(50) NOT NULL,
    department VARCHAR(100),
    position VARCHAR(100),
    phone VARCHAR(20),
    hire_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    permissions JSONB,
    manager_id VARCHAR(100),
    UNIQUE(company_id, employee_code)
);

-- Create face_data table
CREATE TABLE IF NOT EXISTS face_data (
    face_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(employee_id) ON DELETE CASCADE,
    face_embedding BYTEA NOT NULL,
    image_path VARCHAR(500),
    image_name VARCHAR(100),
    image_size INTEGER,
    image_type VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Note: Devices, WorkShifts, AttendanceRecords, and DailyAttendanceSummary tables
-- are not needed for the Identity microservice. These would be handled by
-- separate microservices (Device Service, Attendance Service, etc.)

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_companies_email ON companies(email);
CREATE INDEX IF NOT EXISTS idx_companies_status ON companies(status);
CREATE INDEX IF NOT EXISTS idx_companies_plan ON companies(plan);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_employees_company_id ON employees(company_id);
CREATE INDEX IF NOT EXISTS idx_employees_user_id ON employees(user_id);
CREATE INDEX IF NOT EXISTS idx_employees_employee_code ON employees(employee_code);
CREATE INDEX IF NOT EXISTS idx_employees_manager_id ON employees(manager_id);
CREATE INDEX IF NOT EXISTS idx_face_data_employee_id ON face_data(employee_id);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Triggers for removed tables (devices, work_shifts, daily_attendance_summary) are not needed
