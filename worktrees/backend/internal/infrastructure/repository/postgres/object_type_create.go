package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openfoundry/oms/internal/domain/entity"
	"go.uber.org/zap"
)

// Create creates a new object type
func (r *ObjectTypeRepository) Create(ctx context.Context, objectType *entity.ObjectType) error {
	r.logger.Debug("Creating object type", zap.String("id", objectType.ID.String()))

	propertiesJSON, err := r.marshalProperties(objectType.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	metadataJSON, err := r.marshalMetadata(objectType.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO object_types (
			id, name, display_name, description, category, tags,
			properties, metadata, is_deleted, version,
			created_at, updated_at, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err = r.db.ExecContext(ctx, query,
		objectType.ID,
		objectType.Name,
		objectType.DisplayName,
		objectType.Description,
		objectType.Category,
		objectType.Tags,
		propertiesJSON,
		metadataJSON,
		objectType.IsDeleted,
		objectType.Version,
		objectType.CreatedAt,
		objectType.UpdatedAt,
		objectType.CreatedBy,
		objectType.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to create object type", zap.Error(err))
		return fmt.Errorf("failed to create object type: %w", err)
	}

	r.logger.Info("Object type created successfully", zap.String("id", objectType.ID.String()))
	return nil
}

// marshalProperties marshals properties to JSON
func (r *ObjectTypeRepository) marshalProperties(properties []entity.Property) ([]byte, error) {
	if properties == nil {
		return json.Marshal([]entity.Property{})
	}
	return json.Marshal(properties)
}

// marshalMetadata marshals metadata to JSON
func (r *ObjectTypeRepository) marshalMetadata(metadata map[string]interface{}) ([]byte, error) {
	if metadata == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(metadata)
}