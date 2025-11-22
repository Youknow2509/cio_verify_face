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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS face_profiles;

-- +goose StatementEnd
