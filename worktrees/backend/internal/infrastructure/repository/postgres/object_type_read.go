package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
	"go.uber.org/zap"
)

// GetByID retrieves an object type by ID
func (r *ObjectTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error) {
	r.logger.Debug("Getting object type by ID", zap.String("id", id.String()))

	query := `
		SELECT id, name, display_name, description, category, tags,
			   properties, metadata, is_deleted, version,
			   created_at, updated_at, created_by, updated_by
		FROM object_types 
		WHERE id = $1 AND is_deleted = false`

	var objectType entity.ObjectType
	var propertiesJSON, metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&objectType.ID,
		&objectType.Name,
		&objectType.DisplayName,
		&objectType.Description,
		&objectType.Category,
		&objectType.Tags,
		&propertiesJSON,
		&metadataJSON,
		&objectType.IsDeleted,
		&objectType.Version,
		&objectType.CreatedAt,
		&objectType.UpdatedAt,
		&objectType.CreatedBy,
		&objectType.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrObjectTypeNotFound
		}
		r.logger.Error("Failed to get object type", zap.Error(err))
		return nil, fmt.Errorf("failed to get object type: %w", err)
	}

	if err := r.unmarshalObjectTypeData(&objectType, propertiesJSON, metadataJSON); err != nil {
		return nil, err
	}

	return &objectType, nil
}

// GetByName retrieves an object type by name
func (r *ObjectTypeRepository) GetByName(ctx context.Context, name string) (*entity.ObjectType, error) {
	r.logger.Debug("Getting object type by name", zap.String("name", name))

	query := `
		SELECT id, name, display_name, description, category, tags,
			   properties, metadata, is_deleted, version,
			   created_at, updated_at, created_by, updated_by
		FROM object_types 
		WHERE name = $1 AND is_deleted = false`

	var objectType entity.ObjectType
	var propertiesJSON, metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&objectType.ID,
		&objectType.Name,
		&objectType.DisplayName,
		&objectType.Description,
		&objectType.Category,
		&objectType.Tags,
		&propertiesJSON,
		&metadataJSON,
		&objectType.IsDeleted,
		&objectType.Version,
		&objectType.CreatedAt,
		&objectType.UpdatedAt,
		&objectType.CreatedBy,
		&objectType.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrObjectTypeNotFound
		}
		r.logger.Error("Failed to get object type by name", zap.Error(err))
		return nil, fmt.Errorf("failed to get object type: %w", err)
	}

	if err := r.unmarshalObjectTypeData(&objectType, propertiesJSON, metadataJSON); err != nil {
		return nil, err
	}

	return &objectType, nil
}

// unmarshalObjectTypeData unmarshals JSON data into object type
func (r *ObjectTypeRepository) unmarshalObjectTypeData(
	objectType *entity.ObjectType,
	propertiesJSON, metadataJSON []byte,
) error {
	if err := json.Unmarshal(propertiesJSON, &objectType.Properties); err != nil {
		r.logger.Error("Failed to unmarshal properties", zap.Error(err))
		return fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &objectType.Metadata); err != nil {
		r.logger.Error("Failed to unmarshal metadata", zap.Error(err))
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return nil
}