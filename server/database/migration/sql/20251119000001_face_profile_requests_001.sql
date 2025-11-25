-- +goose Up
-- +goose StatementBegin

-- =================================================================
-- FACE PROFILE UPDATE REQUESTS TABLE
-- =================================================================
-- This table stores requests from employees to update their face profile
-- Partitioned by company_id for better performance and data isolation

CREATE TABLE IF NOT EXISTS face_profile_update_requests (
    request_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    
    -- Request status: 0=pending, 1=approved, 2=rejected, 3=expired, 4=completed
    status INT2 DEFAULT 0 CHECK (status IN (0, 1, 2, 3, 4)),
    
    -- Monthly request tracking (to limit requests per month)
    request_month VARCHAR(7) NOT NULL, -- Format: YYYY-MM
    request_count_in_month INT DEFAULT 1,
    
    -- Update link details (generated when approved)
    update_token VARCHAR(255),
    update_link_expires_at TIMESTAMP WITH TIME ZONE,
    
    -- Approval/Rejection details
    approved_by UUID REFERENCES users(user_id),
    approved_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    
    -- Metadata
    reason TEXT, -- Employee's reason for requesting update
    meta_data JSONB DEFAULT '{}' NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    PRIMARY KEY (request_id, company_id),
    -- Unique constraint must include partitioning column
    UNIQUE (update_token, company_id)
) PARTITION BY LIST (company_id);


-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_fpr_user_id ON face_profile_update_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_fpr_status ON face_profile_update_requests(status);
CREATE INDEX IF NOT EXISTS idx_fpr_request_month ON face_profile_update_requests(request_month);
CREATE INDEX IF NOT EXISTS idx_fpr_update_token ON face_profile_update_requests(update_token) WHERE update_token IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_fpr_user_month ON face_profile_update_requests(user_id, request_month);
CREATE INDEX IF NOT EXISTS idx_fpr_pending_company ON face_profile_update_requests(company_id, status) WHERE status = 0;

-- =================================================================
-- PASSWORD RESET REQUESTS TABLE
-- =================================================================
-- This table stores password reset requests to prevent spam

CREATE TABLE IF NOT EXISTS password_reset_requests (
    request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(company_id) ON DELETE CASCADE,
    
    -- Requested by (manager or admin who initiated the reset)
    requested_by UUID NOT NULL REFERENCES users(user_id),
    
    -- Status: 0=pending, 1=sent, 2=failed
    status INT2 DEFAULT 0 CHECK (status IN (0, 1, 2)),
    
    -- Kafka message tracking
    kafka_message_id VARCHAR(255),
    kafka_sent_at TIMESTAMP WITH TIME ZONE,
    
    -- Metadata for tracking
    meta_data JSONB DEFAULT '{}' NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes for password reset requests
CREATE INDEX IF NOT EXISTS idx_prr_user_id ON password_reset_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_prr_requested_by ON password_reset_requests(requested_by);
CREATE INDEX IF NOT EXISTS idx_prr_company_id ON password_reset_requests(company_id);
CREATE INDEX IF NOT EXISTS idx_prr_created_at ON password_reset_requests(created_at DESC);
-- Index for spam prevention (check recent requests from same manager)
CREATE INDEX IF NOT EXISTS idx_prr_spam_check ON password_reset_requests(requested_by, user_id, created_at DESC);

-- =================================================================
-- CREATE PARTITIONS FOR EXISTING COMPANIES
-- =================================================================
DO $$
DECLARE
    company_rec RECORD;
    partition_name TEXT;
BEGIN
    FOR company_rec IN SELECT company_id FROM companies LOOP
        partition_name := 'fpr_p_' || replace(company_rec.company_id::text, '-', '');
        IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname = partition_name AND relkind = 'r') THEN
            EXECUTE format(
                'CREATE TABLE %I PARTITION OF face_profile_update_requests FOR VALUES IN (%L)',
                partition_name,
                company_rec.company_id
            );
        END IF;
    END LOOP;
END $$;

-- =================================================================
-- FUNCTION TO AUTO-CREATE PARTITION FOR NEW COMPANIES
-- =================================================================
CREATE OR REPLACE FUNCTION create_fpr_partition_for_company()
RETURNS TRIGGER AS $$
DECLARE
    partition_name TEXT;
BEGIN
    partition_name := 'fpr_p_' || replace(NEW.company_id::text, '-', '');
    IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname = partition_name AND relkind = 'r') THEN
        EXECUTE format(
            'CREATE TABLE %I PARTITION OF face_profile_update_requests FOR VALUES IN (%L)',
            partition_name,
            NEW.company_id
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-create partition when new company is added
DROP TRIGGER IF EXISTS trg_create_fpr_partition ON companies;
CREATE TRIGGER trg_create_fpr_partition
    AFTER INSERT ON companies
    FOR EACH ROW
    EXECUTE FUNCTION create_fpr_partition_for_company();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger and function
DROP TRIGGER IF EXISTS trg_create_fpr_partition ON companies;
DROP FUNCTION IF EXISTS create_fpr_partition_for_company();

-- Drop tables (partitions will be dropped automatically)
DROP TABLE IF EXISTS face_profile_update_requests CASCADE;
DROP TABLE IF EXISTS password_reset_requests CASCADE;

-- +goose StatementEnd
