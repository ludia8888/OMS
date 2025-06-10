package postgres

import (
	"context"
	"fmt"

	"github.com/openfoundry/oms/internal/domain/entity"
	"go.uber.org/zap"
)

// Update updates an existing object type
func (r *ObjectTypeRepository) Update(ctx context.Context, objectType *entity.ObjectType) error {
	r.logger.Debug("Updating object type", zap.String("id", objectType.ID.String()))

	propertiesJSON, err := r.marshalProperties(objectType.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	metadataJSON, err := r.marshalMetadata(objectType.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		UPDATE object_types 
		SET display_name = $2, description = $3, category = $4, tags = $5,
			properties = $6, metadata = $7, version = $8,
			updated_at = $9, updated_by = $10
		WHERE id = $1 AND is_deleted = false`

	result, err := r.db.ExecContext(ctx, query,
		objectType.ID,
		objectType.DisplayName,
		objectType.Description,
		objectType.Category,
		objectType.Tags,
		propertiesJSON,
		metadataJSON,
		objectType.Version,
		objectType.UpdatedAt,
		objectType.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to update object type", zap.Error(err))
		return fmt.Errorf("failed to update object type: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrObjectTypeNotFound
	}

	r.logger.Info("Object type updated successfully", zap.String("id", objectType.ID.String()))
	return nil
}