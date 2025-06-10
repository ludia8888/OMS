-- Create link_types table
CREATE TABLE IF NOT EXISTS link_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    source_object_type_id UUID NOT NULL REFERENCES object_types(id),
    target_object_type_id UUID NOT NULL REFERENCES object_types(id),
    cardinality VARCHAR(32) NOT NULL CHECK (cardinality IN ('ONE_TO_ONE', 'ONE_TO_MANY', 'MANY_TO_MANY')),
    description TEXT,
    properties JSONB DEFAULT '[]'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    version INTEGER NOT NULL DEFAULT 1,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL,
    
    CONSTRAINT link_type_name_format CHECK (name ~ '^[a-zA-Z][a-zA-Z0-9_]*$')
);

-- Create indexes
CREATE INDEX idx_link_types_name ON link_types(name) WHERE is_deleted = FALSE;
CREATE INDEX idx_link_types_source ON link_types(source_object_type_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_link_types_target ON link_types(target_object_type_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_link_types_source_target ON link_types(source_object_type_id, target_object_type_id) WHERE is_deleted = FALSE;

-- Create link_type_versions table for version history
CREATE TABLE IF NOT EXISTS link_type_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_type_id UUID NOT NULL REFERENCES link_types(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    snapshot JSONB NOT NULL,
    change_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    
    UNIQUE(link_type_id, version)
);

-- Create index for version history queries
CREATE INDEX idx_link_type_versions_link_type_id ON link_type_versions(link_type_id);
CREATE INDEX idx_link_type_versions_created_at ON link_type_versions(created_at DESC);