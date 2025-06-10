package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
)

// LinkTypeRepository defines the interface for link type persistence
type LinkTypeRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, linkType *entity.LinkType) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.LinkType, error)
	GetByName(ctx context.Context, name string) (*entity.LinkType, error)
	Update(ctx context.Context, linkType *entity.LinkType) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	List(ctx context.Context, filter LinkTypeFilter) ([]*entity.LinkType, error)
	Count(ctx context.Context, filter LinkTypeFilter) (int64, error)

	// Relationship queries
	GetBySourceObjectType(ctx context.Context, objectTypeID uuid.UUID) ([]*entity.LinkType, error)
	GetByTargetObjectType(ctx context.Context, objectTypeID uuid.UUID) ([]*entity.LinkType, error)
	GetByObjectTypes(ctx context.Context, sourceID, targetID uuid.UUID) ([]*entity.LinkType, error)

	// Validation
	CheckCircularReference(ctx context.Context, sourceID, targetID uuid.UUID) (bool, error)
}

// LinkTypeFilter represents filtering options for link types
type LinkTypeFilter struct {
	SourceObjectTypeID *uuid.UUID
	TargetObjectTypeID *uuid.UUID
	Cardinality       *entity.Cardinality
	IsDeleted         *bool
	PageSize          int
	PageCursor        string
	SortBy            string
	SortOrder         string
}