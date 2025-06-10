-- Create object_types table
CREATE TABLE IF NOT EXISTS object_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(64),
    tags JSONB DEFAULT '[]'::jsonb,
    properties JSONB NOT NULL DEFAULT '[]'::jsonb,
    base_datasets JSONB DEFAULT '[]'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    version INTEGER NOT NULL DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    
    CONSTRAINT object_type_name_format CHECK (name ~ '^[a-zA-Z][a-zA-Z0-9_]*$'),
    CONSTRAINT object_type_name_length CHECK (char_length(name) <= 64)
);

-- Create indexes
CREATE INDEX idx_object_types_name ON object_types(name) WHERE is_deleted = FALSE;
CREATE INDEX idx_object_types_category ON object_types(category) WHERE is_deleted = FALSE;
CREATE INDEX idx_object_types_tags ON object_types USING GIN (tags) WHERE is_deleted = FALSE;
CREATE INDEX idx_object_types_created_at ON object_types(created_at DESC) WHERE is_deleted = FALSE;

-- Full-text search index
CREATE INDEX idx_object_types_search ON object_types 
USING GIN (to_tsvector('english', name || ' ' || display_name || ' ' || COALESCE(description, ''))) 
WHERE is_deleted = FALSE;

-- Create object_type_versions table for version history
CREATE TABLE IF NOT EXISTS object_type_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object_type_id UUID NOT NULL REFERENCES object_types(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    snapshot JSONB NOT NULL,
    change_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    
    UNIQUE(object_type_id, version)
);

-- Create index for version history queries
CREATE INDEX idx_object_type_versions_object_type_id ON object_type_versions(object_type_id);
CREATE INDEX idx_object_type_versions_created_at ON object_type_versions(created_at DESC);