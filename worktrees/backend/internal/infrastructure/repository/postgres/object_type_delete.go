package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
	"go.uber.org/zap"
)

// Delete soft deletes an object type
func (r *ObjectTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Deleting object type", zap.String("id", id.String()))

	query := `
		UPDATE object_types 
		SET is_deleted = true, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_deleted = false`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete object type", zap.Error(err))
		return fmt.Errorf("failed to delete object type: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrObjectTypeNotFound
	}

	r.logger.Info("Object type deleted successfully", zap.String("id", id.String()))
	return nil
}