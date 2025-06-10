package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
)

// ObjectTypeRepository defines the interface for object type persistence
type ObjectTypeRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, objectType *entity.ObjectType) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error)
	GetByName(ctx context.Context, name string) (*entity.ObjectType, error)
	Update(ctx context.Context, objectType *entity.ObjectType) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	List(ctx context.Context, filter ObjectTypeFilter) ([]*entity.ObjectType, error)
	Count(ctx context.Context, filter ObjectTypeFilter) (int64, error)
	Search(ctx context.Context, query string, limit int) ([]*entity.ObjectType, error)

	// Version management
	GetVersion(ctx context.Context, id uuid.UUID, version int) (*entity.ObjectType, error)
	ListVersions(ctx context.Context, id uuid.UUID) ([]*ObjectTypeVersion, error)
	CompareVersions(ctx context.Context, id uuid.UUID, v1, v2 int) (*VersionDiff, error)

	// Batch operations
	BatchCreate(ctx context.Context, objectTypes []*entity.ObjectType) error
	BatchUpdate(ctx context.Context, objectTypes []*entity.ObjectType) error
}

// ObjectTypeFilter represents filtering options for object types
type ObjectTypeFilter struct {
	Category      *string
	Tags          []string
	IsDeleted     *bool
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
	PageSize      int
	PageCursor    string // Cursor-based pagination
	SortBy        string
	SortOrder     string // "asc" or "desc"
}

// ObjectTypeVersion represents a historical version of an object type
type ObjectTypeVersion struct {
	ID               uuid.UUID            `json:"id"`
	ObjectTypeID     uuid.UUID            `json:"objectTypeId"`
	Version          int                  `json:"version"`
	Snapshot         entity.ObjectType    `json:"snapshot"`
	ChangeDescription string              `json:"changeDescription,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
	CreatedBy        string              `json:"createdBy"`
}

// VersionDiff represents the difference between two versions
type VersionDiff struct {
	ObjectTypeID uuid.UUID      `json:"objectTypeId"`
	Version1     int            `json:"version1"`
	Version2     int            `json:"version2"`
	Changes      []FieldChange  `json:"changes"`
}

// FieldChange represents a change in a field
type FieldChange struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"oldValue"`
	NewValue interface{} `json:"newValue"`
	Type     ChangeType  `json:"type"`
}

// ChangeType represents the type of change
type ChangeType string

const (
	ChangeTypeAdded    ChangeType = "added"
	ChangeTypeRemoved  ChangeType = "removed"
	ChangeTypeModified ChangeType = "modified"
)

// PageCursor represents pagination cursor information
type PageCursor struct {
	Timestamp time.Time
	ID        uuid.UUID
}