-- +goose Up
-- +goose StatementBegin

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Add vector column to face_profiles table if it doesn't exist
-- This migration is safe to run multiple times
DO $$ 
BEGIN
    -- Check if the embedding column exists and is of type FLOAT4[]
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'face_profiles' 
          AND column_name = 'embedding'
          AND data_type = 'ARRAY'
    ) THEN
        -- The column exists as FLOAT4[], we need to migrate data to vector type
        
        -- First, create a temporary column for the vector data
        ALTER TABLE face_profiles ADD COLUMN IF NOT EXISTS embedding_vector vector(512);
        
        -- Convert existing FLOAT4[] data to vector type
        UPDATE face_profiles 
        SET embedding_vector = embedding::text::vector
        WHERE embedding IS NOT NULL AND embedding_vector IS NULL;
        
        -- Drop the old FLOAT4[] column
        ALTER TABLE face_profiles DROP COLUMN IF EXISTS embedding;
        
        -- Rename the new vector column to embedding
        ALTER TABLE face_profiles RENAME COLUMN embedding_vector TO embedding;
        
        -- Ensure the column is NOT NULL
        ALTER TABLE face_profiles ALTER COLUMN embedding SET NOT NULL;
        
        RAISE NOTICE 'Successfully migrated embedding column from FLOAT4[] to vector type';
    ELSIF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'face_profiles' 
          AND column_name = 'embedding'
    ) THEN
        -- Column doesn't exist at all, create it
        ALTER TABLE face_profiles ADD COLUMN embedding vector(512) NOT NULL;
        RAISE NOTICE 'Created new embedding column with vector type';
    END IF;
END $$;

-- Create an index for vector similarity search using cosine distance
-- Drop existing index if it exists
DROP INDEX IF EXISTS idx_face_profiles_embedding_cosine;

-- Create HNSW index for fast approximate nearest neighbor search
-- HNSW (Hierarchical Navigable Small World) is best for high-dimensional vectors
-- m=16 is a good default (higher = better recall but more memory)
-- ef_construction=64 is a good default (higher = better index quality but slower build)
CREATE INDEX idx_face_profiles_embedding_cosine 
ON face_profiles 
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- Alternative: IVFFlat index (uncomment if you prefer this over HNSW)
-- IVFFlat is faster to build but slightly slower for queries
-- DROP INDEX IF EXISTS idx_face_profiles_embedding_cosine;
-- CREATE INDEX idx_face_profiles_embedding_cosine 
-- ON face_profiles 
-- USING ivfflat (embedding vector_cosine_ops)
-- WITH (lists = 100);

-- Add index on indexed flag for faster filtering
CREATE INDEX IF NOT EXISTS idx_face_profiles_indexed 
ON face_profiles(indexed) 
WHERE indexed = true AND deleted_at IS NULL;

-- Update statistics for query planner
ANALYZE face_profiles;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove vector-specific indexes
DROP INDEX IF EXISTS idx_face_profiles_embedding_cosine;
DROP INDEX IF EXISTS idx_face_profiles_indexed;

-- Convert vector column back to FLOAT4[] array
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'face_profiles' 
          AND column_name = 'embedding'
          AND udt_name = 'vector'
    ) THEN
        -- Create temporary column for array data
        ALTER TABLE face_profiles ADD COLUMN IF NOT EXISTS embedding_array FLOAT4[];
        
        -- Convert vector data back to array
        -- Note: This is a lossy conversion as vector is stored more efficiently
        UPDATE face_profiles 
        SET embedding_array = (embedding::text::float[])::FLOAT4[]
        WHERE embedding IS NOT NULL;
        
        -- Drop vector column
        ALTER TABLE face_profiles DROP COLUMN embedding;
        
        -- Rename array column back to embedding
        ALTER TABLE face_profiles RENAME COLUMN embedding_array TO embedding;
        
        -- Ensure NOT NULL constraint
        ALTER TABLE face_profiles ALTER COLUMN embedding SET NOT NULL;
    END IF;
END $$;

-- Note: We don't drop the vector extension in case other tables use it
-- DROP EXTENSION IF EXISTS vector;

-- +goose StatementEnd
