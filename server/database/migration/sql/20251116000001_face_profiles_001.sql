-- +goose Up
-- +goose StatementBegin

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;
-- Create face_profiles table
CREATE TABLE face_profiles (
   profile_id UUID NOT NULL,
   user_id UUID NOT NULL,
   company_id UUID,
   embedding vector(512) NOT NULL,
   embedding_version VARCHAR(50) NOT NULL,
   enroll_image_path TEXT,
   is_primary BOOLEAN DEFAULT false NOT NULL,
   quality_score FLOAT,
   meta_data JSONB DEFAULT '{}' NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
   deleted_at TIMESTAMP WITH TIME ZONE,
   indexed BOOLEAN DEFAULT false NOT NULL,
   index_version INTEGER DEFAULT 0 NOT NULL,
   PRIMARY KEY (profile_id, company_id)
) PARTITION BY LIST (company_id);

--  Create partitions for each company
--  ...........

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_face_profiles_user_id ON face_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_face_profiles_embedding_version ON face_profiles(embedding_version);
CREATE INDEX IF NOT EXISTS idx_face_profiles_deleted_at ON face_profiles(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_face_profiles_is_primary ON face_profiles(is_primary) WHERE is_primary = TRUE;

-- Create face_audit_logs table
CREATE TABLE IF NOT EXISTS face_audit_logs (
    log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID,
    user_id UUID,
    
    -- Operation details
    operation VARCHAR(50) NOT NULL,  -- enroll, verify, update, delete
    status VARCHAR(20) NOT NULL,  -- success, failed, duplicate, no_match
    
    -- Request metadata
    device_id VARCHAR(100),
    ip_address VARCHAR(50),
    
    -- Results
    similarity_score FLOAT,
    liveness_score FLOAT,
    quality_score FLOAT,
    
    -- Additional data
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    error_message TEXT,
    
    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes for audit logs
CREATE INDEX IF NOT EXISTS idx_face_audit_logs_profile_id ON face_audit_logs(profile_id);
CREATE INDEX IF NOT EXISTS idx_face_audit_logs_user_id ON face_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_face_audit_logs_operation ON face_audit_logs(operation);
CREATE INDEX IF NOT EXISTS idx_face_audit_logs_created_at ON face_audit_logs(created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS face_profiles_updated_at ON face_profiles;
DROP TABLE IF EXISTS face_audit_logs;
DROP TABLE IF EXISTS face_profiles;

-- +goose StatementEnd
