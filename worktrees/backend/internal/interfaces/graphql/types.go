package graphql

import (
	"time"

	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
	"github.com/openfoundry/oms/internal/domain/repository"
)

// Query types

type ObjectTypeFilter struct {
	Category      *string    `json:"category"`
	Tags          []string   `json:"tags"`
	IsDeleted     *bool      `json:"isDeleted"`
	CreatedAfter  *time.Time `json:"createdAfter"`
	CreatedBefore *time.Time `json:"createdBefore"`
	UpdatedAfter  *time.Time `json:"updatedAfter"`
	UpdatedBefore *time.Time `json:"updatedBefore"`
}

type LinkTypeFilter struct {
	SourceObjectTypeID *uuid.UUID         `json:"sourceObjectTypeId"`
	TargetObjectTypeID *uuid.UUID         `json:"targetObjectTypeId"`
	Cardinality        *entity.Cardinality `json:"cardinality"`
	IsDeleted          *bool              `json:"isDeleted"`
}

type PaginationInput struct {
	PageSize  *int       `json:"pageSize"`
	Cursor    *string    `json:"cursor"`
	SortBy    *string    `json:"sortBy"`
	SortOrder *SortOrder `json:"sortOrder"`
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

type ObjectTypeConnection struct {
	Edges      []*ObjectTypeEdge `json:"edges"`
	PageInfo   *PageInfo         `json:"pageInfo"`
	TotalCount int               `json:"totalCount"`
}

type ObjectTypeEdge struct {
	Node   *entity.ObjectType `json:"node"`
	Cursor string             `json:"cursor"`
}

type LinkTypeConnection struct {
	Edges      []*LinkTypeEdge `json:"edges"`
	PageInfo   *PageInfo       `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

type LinkTypeEdge struct {
	Node   *entity.LinkType `json:"node"`
	Cursor string           `json:"cursor"`
}

// Mutation inputs

type CreateObjectTypeInput struct {
	Name        string                 `json:"name"`
	DisplayName string                 `json:"displayName"`
	Description *string                `json:"description"`
	Category    *string                `json:"category"`
	Properties  []*CreatePropertyInput `json:"properties"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UpdateObjectTypeInput struct {
	DisplayName *string                `json:"displayName"`
	Description *string                `json:"description"`
	Category    *string                `json:"category"`
	Properties  []*UpdatePropertyInput `json:"properties"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreatePropertyInput struct {
	Name            string                 `json:"name"`
	DisplayName     string                 `json:"displayName"`
	DataType        entity.DataType        `json:"dataType"`
	IsRequired      *bool                  `json:"isRequired"`
	IsArray         *bool                  `json:"isArray"`
	IsUnique        *bool                  `json:"isUnique"`
	DefaultValue    interface{}            `json:"defaultValue"`
	ValidationRules *ValidationRulesInput  `json:"validationRules"`
	Metadata        map[string]interface{} `json:"metadata"`
	Order           *int                   `json:"order"`
}

type UpdatePropertyInput struct {
	ID              *uuid.UUID             `json:"id"`
	Name            *string                `json:"name"`
	DisplayName     *string                `json:"displayName"`
	DataType        *entity.DataType       `json:"dataType"`
	IsRequired      *bool                  `json:"isRequired"`
	IsArray         *bool                  `json:"isArray"`
	IsUnique        *bool                  `json:"isUnique"`
	DefaultValue    interface{}            `json:"defaultValue"`
	ValidationRules *ValidationRulesInput  `json:"validationRules"`
	Metadata        map[string]interface{} `json:"metadata"`
	Order           *int                   `json:"order"`
}

type ValidationRulesInput struct {
	Pattern     *string                `json:"pattern"`
	MinLength   *int                   `json:"minLength"`
	MaxLength   *int                   `json:"maxLength"`
	MinValue    *float64               `json:"minValue"`
	MaxValue    *float64               `json:"maxValue"`
	EnumValues  []string               `json:"enumValues"`
	CustomRules map[string]interface{} `json:"customRules"`
}

type CreateLinkTypeInput struct {
	Name               string                 `json:"name"`
	DisplayName        string                 `json:"displayName"`
	InverseDisplayName *string                `json:"inverseDisplayName"`
	Description        *string                `json:"description"`
	SourceObjectTypeID uuid.UUID              `json:"sourceObjectTypeId"`
	TargetObjectTypeID uuid.UUID              `json:"targetObjectTypeId"`
	Cardinality        entity.Cardinality     `json:"cardinality"`
	Properties         []*CreatePropertyInput `json:"properties"`
	Constraints        *LinkConstraintsInput  `json:"constraints"`
	Metadata           map[string]interface{} `json:"metadata"`
}

type UpdateLinkTypeInput struct {
	DisplayName        *string                `json:"displayName"`
	InverseDisplayName *string                `json:"inverseDisplayName"`
	Description        *string                `json:"description"`
	Cardinality        *entity.Cardinality    `json:"cardinality"`
	Properties         []*UpdatePropertyInput `json:"properties"`
	Constraints        *LinkConstraintsInput  `json:"constraints"`
	Metadata           map[string]interface{} `json:"metadata"`
}

type LinkConstraintsInput struct {
	IsRequired      *bool                  `json:"isRequired"`
	CascadeDelete   *bool                  `json:"cascadeDelete"`
	PreventDelete   *bool                  `json:"preventDelete"`
	UniquePerSource *bool                  `json:"uniquePerSource"`
	UniquePerTarget *bool                  `json:"uniquePerTarget"`
	ValidationRules map[string]interface{} `json:"validationRules"`
}

// Resolver interfaces

type QueryResolver interface {
	ObjectType(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error)
	ObjectTypes(ctx context.Context, filter *ObjectTypeFilter, pagination *PaginationInput) (*ObjectTypeConnection, error)
	SearchObjectTypes(ctx context.Context, query string, limit *int) ([]*entity.ObjectType, error)
	LinkType(ctx context.Context, id uuid.UUID) (*entity.LinkType, error)
	LinkTypes(ctx context.Context, filter *LinkTypeFilter, pagination *PaginationInput) (*LinkTypeConnection, error)
	LinkTypesByObjectTypes(ctx context.Context, sourceID uuid.UUID, targetID uuid.UUID) ([]*entity.LinkType, error)
	ObjectTypeVersion(ctx context.Context, objectTypeID uuid.UUID, version int) (*repository.ObjectTypeVersion, error)
	ObjectTypeVersions(ctx context.Context, objectTypeID uuid.UUID) ([]*repository.ObjectTypeVersion, error)
	CompareObjectTypeVersions(ctx context.Context, objectTypeID uuid.UUID, v1 int, v2 int) (*repository.VersionDiff, error)
}

type MutationResolver interface {
	CreateObjectType(ctx context.Context, input CreateObjectTypeInput) (*entity.ObjectType, error)
	UpdateObjectType(ctx context.Context, id uuid.UUID, input UpdateObjectTypeInput) (*entity.ObjectType, error)
	DeleteObjectType(ctx context.Context, id uuid.UUID) (bool, error)
	CreateLinkType(ctx context.Context, input CreateLinkTypeInput) (*entity.LinkType, error)
	UpdateLinkType(ctx context.Context, id uuid.UUID, input UpdateLinkTypeInput) (*entity.LinkType, error)
	DeleteLinkType(ctx context.Context, id uuid.UUID) (bool, error)
}

type ObjectTypeResolver interface {
	OutgoingLinkTypes(ctx context.Context, obj *entity.ObjectType) ([]*entity.LinkType, error)
	IncomingLinkTypes(ctx context.Context, obj *entity.ObjectType) ([]*entity.LinkType, error)
}

type LinkTypeResolver interface {
	SourceObjectType(ctx context.Context, obj *entity.LinkType) (*entity.ObjectType, error)
	TargetObjectType(ctx context.Context, obj *entity.LinkType) (*entity.ObjectType, error)
}