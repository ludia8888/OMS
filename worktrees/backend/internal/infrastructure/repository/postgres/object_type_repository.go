package postgres

import (
	"database/sql"

	"github.com/openfoundry/oms/internal/domain/repository"
	"go.uber.org/zap"
)

// ObjectTypeRepository implements repository.ObjectTypeRepository using PostgreSQL
type ObjectTypeRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewObjectTypeRepository creates a new PostgreSQL object type repository
func NewObjectTypeRepository(db *sql.DB, logger *zap.Logger) repository.ObjectTypeRepository {
	return &ObjectTypeRepository{
		db:     db,
		logger: logger,
	}
}