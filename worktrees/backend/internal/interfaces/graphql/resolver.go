package graphql

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
	"github.com/openfoundry/oms/internal/domain/repository"
	"github.com/openfoundry/oms/internal/domain/service"
	"github.com/openfoundry/oms/internal/interfaces/rest/middleware"
	"go.uber.org/zap"
)

// Resolver is the root resolver
type Resolver struct {
	objectTypeService *service.ObjectTypeService
	linkTypeService   *service.LinkTypeService
	logger            *zap.Logger
}

// NewResolver creates a new GraphQL resolver
func NewResolver(objectTypeService *service.ObjectTypeService, linkTypeService *service.LinkTypeService, logger *zap.Logger) *Resolver {
	return &Resolver{
		objectTypeService: objectTypeService,
		linkTypeService:   linkTypeService,
		logger:            logger,
	}
}

// Query returns the query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Mutation returns the mutation resolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// ObjectType returns the object type resolver
func (r *Resolver) ObjectType() ObjectTypeResolver {
	return &objectTypeResolver{r}
}

// LinkType returns the link type resolver
func (r *Resolver) LinkType() LinkTypeResolver {
	return &linkTypeResolver{r}
}

type queryResolver struct{ *Resolver }

// ObjectType retrieves an object type by ID
func (r *queryResolver) ObjectType(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error) {
	return r.objectTypeService.GetByID(ctx, id)
}

// ObjectTypes lists object types with filtering and pagination
func (r *queryResolver) ObjectTypes(ctx context.Context, filter *ObjectTypeFilter, pagination *PaginationInput) (*ObjectTypeConnection, error) {
	// Convert GraphQL filter to repository filter
	repoFilter := repository.ObjectTypeFilter{}
	
	if filter != nil {
		repoFilter.Category = filter.Category
		repoFilter.Tags = filter.Tags
		repoFilter.IsDeleted = filter.IsDeleted
		repoFilter.CreatedAfter = filter.CreatedAfter
		repoFilter.CreatedBefore = filter.CreatedBefore
		repoFilter.UpdatedAfter = filter.UpdatedAfter
		repoFilter.UpdatedBefore = filter.UpdatedBefore
	}

	if pagination != nil {
		repoFilter.PageSize = getPageSize(pagination.PageSize)
		repoFilter.PageCursor = getString(pagination.Cursor)
		repoFilter.SortBy = getString(pagination.SortBy)
		repoFilter.SortOrder = getSortOrder(pagination.SortOrder)
	} else {
		repoFilter.PageSize = 20
	}

	// Get object types
	objectTypes, err := r.objectTypeService.List(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	// Build connection response
	edges := make([]*ObjectTypeEdge, len(objectTypes))
	for i, ot := range objectTypes {
		edges[i] = &ObjectTypeEdge{
			Node:   ot,
			Cursor: encodeCursor(ot.CreatedAt, ot.ID),
		}
	}

	pageInfo := &PageInfo{
		HasNextPage: len(objectTypes) == repoFilter.PageSize,
	}

	if len(edges) > 0 {
		pageInfo.StartCursor = &edges[0].Cursor
		pageInfo.EndCursor = &edges[len(edges)-1].Cursor
	}

	// Get total count
	count, err := r.objectTypeService.Count(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	return &ObjectTypeConnection{
		Edges:      edges,
		PageInfo:   pageInfo,
		TotalCount: int(count),
	}, nil
}

// SearchObjectTypes searches for object types
func (r *queryResolver) SearchObjectTypes(ctx context.Context, query string, limit *int) ([]*entity.ObjectType, error) {
	searchLimit := 10
	if limit != nil && *limit > 0 && *limit <= 50 {
		searchLimit = *limit
	}
	return r.objectTypeService.Search(ctx, query, searchLimit)
}

// LinkType retrieves a link type by ID
func (r *queryResolver) LinkType(ctx context.Context, id uuid.UUID) (*entity.LinkType, error) {
	return r.linkTypeService.GetByID(ctx, id)
}

// LinkTypes lists link types with filtering and pagination
func (r *queryResolver) LinkTypes(ctx context.Context, filter *LinkTypeFilter, pagination *PaginationInput) (*LinkTypeConnection, error) {
	// Convert GraphQL filter to repository filter
	repoFilter := repository.LinkTypeFilter{}
	
	if filter != nil {
		if filter.SourceObjectTypeID != nil {
			repoFilter.SourceObjectTypeID = filter.SourceObjectTypeID
		}
		if filter.TargetObjectTypeID != nil {
			repoFilter.TargetObjectTypeID = filter.TargetObjectTypeID
		}
		if filter.Cardinality != nil {
			card := entity.Cardinality(*filter.Cardinality)
			repoFilter.Cardinality = &card
		}
		repoFilter.IsDeleted = filter.IsDeleted
	}

	if pagination != nil {
		repoFilter.PageSize = getPageSize(pagination.PageSize)
		repoFilter.PageCursor = getString(pagination.Cursor)
		repoFilter.SortBy = getString(pagination.SortBy)
		repoFilter.SortOrder = getSortOrder(pagination.SortOrder)
	} else {
		repoFilter.PageSize = 20
	}

	// Get link types
	linkTypes, err := r.linkTypeService.List(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	// Build connection response
	edges := make([]*LinkTypeEdge, len(linkTypes))
	for i, lt := range linkTypes {
		edges[i] = &LinkTypeEdge{
			Node:   lt,
			Cursor: encodeCursor(lt.CreatedAt, lt.ID),
		}
	}

	pageInfo := &PageInfo{
		HasNextPage: len(linkTypes) == repoFilter.PageSize,
	}

	if len(edges) > 0 {
		pageInfo.StartCursor = &edges[0].Cursor
		pageInfo.EndCursor = &edges[len(edges)-1].Cursor
	}

	// Get total count
	count, err := r.linkTypeService.Count(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	return &LinkTypeConnection{
		Edges:      edges,
		PageInfo:   pageInfo,
		TotalCount: int(count),
	}, nil
}

// LinkTypesByObjectTypes retrieves link types between two object types
func (r *queryResolver) LinkTypesByObjectTypes(ctx context.Context, sourceID uuid.UUID, targetID uuid.UUID) ([]*entity.LinkType, error) {
	return r.linkTypeService.GetByObjectTypes(ctx, sourceID, targetID)
}

// ObjectTypeVersion retrieves a specific version of an object type
func (r *queryResolver) ObjectTypeVersion(ctx context.Context, objectTypeID uuid.UUID, version int) (*repository.ObjectTypeVersion, error) {
	return r.objectTypeService.GetVersion(ctx, objectTypeID, version)
}

// ObjectTypeVersions lists all versions of an object type
func (r *queryResolver) ObjectTypeVersions(ctx context.Context, objectTypeID uuid.UUID) ([]*repository.ObjectTypeVersion, error) {
	return r.objectTypeService.ListVersions(ctx, objectTypeID)
}

// CompareObjectTypeVersions compares two versions of an object type
func (r *queryResolver) CompareObjectTypeVersions(ctx context.Context, objectTypeID uuid.UUID, v1 int, v2 int) (*repository.VersionDiff, error) {
	return r.objectTypeService.CompareVersions(ctx, objectTypeID, v1, v2)
}

type mutationResolver struct{ *Resolver }

// CreateObjectType creates a new object type
func (r *mutationResolver) CreateObjectType(ctx context.Context, input CreateObjectTypeInput) (*entity.ObjectType, error) {
	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrUnauthorized
	}

	// Convert input to service input
	serviceInput := service.CreateObjectTypeInput{
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Description: input.Description,
		Category:    input.Category,
		Tags:        input.Tags,
		Metadata:    input.Metadata,
	}

	// Convert properties
	if input.Properties != nil {
		serviceInput.Properties = make([]entity.Property, len(input.Properties))
		for i, prop := range input.Properties {
			serviceInput.Properties[i] = convertPropertyInput(prop)
		}
	}

	return r.objectTypeService.CreateObjectType(ctx, serviceInput, userID)
}

// UpdateObjectType updates an existing object type
func (r *mutationResolver) UpdateObjectType(ctx context.Context, id uuid.UUID, input UpdateObjectTypeInput) (*entity.ObjectType, error) {
	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrUnauthorized
	}

	// Convert input to service input
	serviceInput := service.UpdateObjectTypeInput{
		DisplayName: input.DisplayName,
		Description: input.Description,
		Category:    input.Category,
		Tags:        input.Tags,
		Metadata:    input.Metadata,
	}

	// Convert properties
	if input.Properties != nil {
		properties := make([]entity.Property, len(input.Properties))
		for i, prop := range input.Properties {
			properties[i] = convertUpdatePropertyInput(prop)
		}
		serviceInput.Properties = &properties
	}

	return r.objectTypeService.UpdateObjectType(ctx, id, serviceInput, userID)
}

// DeleteObjectType deletes an object type
func (r *mutationResolver) DeleteObjectType(ctx context.Context, id uuid.UUID) (bool, error) {
	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return false, ErrUnauthorized
	}

	// Check admin role
	if !hasAdminRole(ctx) {
		return false, ErrForbidden
	}

	err := r.objectTypeService.DeleteObjectType(ctx, id, userID)
	return err == nil, err
}

// CreateLinkType creates a new link type
func (r *mutationResolver) CreateLinkType(ctx context.Context, input CreateLinkTypeInput) (*entity.LinkType, error) {
	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrUnauthorized
	}

	// Convert input to service input
	serviceInput := service.CreateLinkTypeInput{
		Name:               input.Name,
		DisplayName:        input.DisplayName,
		InverseDisplayName: input.InverseDisplayName,
		Description:        input.Description,
		SourceObjectTypeID: input.SourceObjectTypeID,
		TargetObjectTypeID: input.TargetObjectTypeID,
		Cardinality:        entity.Cardinality(input.Cardinality),
		Metadata:           input.Metadata,
	}

	// Convert properties
	if input.Properties != nil {
		serviceInput.Properties = make([]entity.Property, len(input.Properties))
		for i, prop := range input.Properties {
			serviceInput.Properties[i] = convertPropertyInput(prop)
		}
	}

	// Convert constraints
	if input.Constraints != nil {
		serviceInput.Constraints = convertLinkConstraintsInput(input.Constraints)
	}

	return r.linkTypeService.CreateLinkType(ctx, serviceInput, userID)
}

// UpdateLinkType updates an existing link type
func (r *mutationResolver) UpdateLinkType(ctx context.Context, id uuid.UUID, input UpdateLinkTypeInput) (*entity.LinkType, error) {
	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrUnauthorized
	}

	// Convert input to service input
	serviceInput := service.UpdateLinkTypeInput{
		DisplayName:        input.DisplayName,
		InverseDisplayName: input.InverseDisplayName,
		Description:        input.Description,
		Metadata:           input.Metadata,
	}

	if input.Cardinality != nil {
		card := entity.Cardinality(*input.Cardinality)
		serviceInput.Cardinality = &card
	}

	// Convert properties
	if input.Properties != nil {
		properties := make([]entity.Property, len(input.Properties))
		for i, prop := range input.Properties {
			properties[i] = convertUpdatePropertyInput(prop)
		}
		serviceInput.Properties = &properties
	}

	// Convert constraints
	if input.Constraints != nil {
		constraints := convertLinkConstraintsInput(input.Constraints)
		serviceInput.Constraints = &constraints
	}

	return r.linkTypeService.UpdateLinkType(ctx, id, serviceInput, userID)
}

// DeleteLinkType deletes a link type
func (r *mutationResolver) DeleteLinkType(ctx context.Context, id uuid.UUID) (bool, error) {
	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return false, ErrUnauthorized
	}

	// Check admin role
	if !hasAdminRole(ctx) {
		return false, ErrForbidden
	}

	err := r.linkTypeService.DeleteLinkType(ctx, id, userID)
	return err == nil, err
}

type objectTypeResolver struct{ *Resolver }

// OutgoingLinkTypes resolves outgoing link types for an object type
func (r *objectTypeResolver) OutgoingLinkTypes(ctx context.Context, obj *entity.ObjectType) ([]*entity.LinkType, error) {
	return r.linkTypeService.GetBySourceObjectType(ctx, obj.ID)
}

// IncomingLinkTypes resolves incoming link types for an object type
func (r *objectTypeResolver) IncomingLinkTypes(ctx context.Context, obj *entity.ObjectType) ([]*entity.LinkType, error) {
	return r.linkTypeService.GetByTargetObjectType(ctx, obj.ID)
}

type linkTypeResolver struct{ *Resolver }

// SourceObjectType resolves the source object type for a link type
func (r *linkTypeResolver) SourceObjectType(ctx context.Context, obj *entity.LinkType) (*entity.ObjectType, error) {
	return r.objectTypeService.GetByID(ctx, obj.SourceObjectTypeID)
}

// TargetObjectType resolves the target object type for a link type
func (r *linkTypeResolver) TargetObjectType(ctx context.Context, obj *entity.LinkType) (*entity.ObjectType, error) {
	return r.objectTypeService.GetByID(ctx, obj.TargetObjectTypeID)
}

// Helper functions

func getUserIDFromContext(ctx context.Context) string {
	if ginCtx := graphql.GetFieldContext(ctx).Args["ginContext"]; ginCtx != nil {
		if gc, ok := ginCtx.(*gin.Context); ok {
			return middleware.GetUserID(gc)
		}
	}
	return ""
}

func hasAdminRole(ctx context.Context) bool {
	if ginCtx := graphql.GetFieldContext(ctx).Args["ginContext"]; ginCtx != nil {
		if gc, ok := ginCtx.(*gin.Context); ok {
			return middleware.HasRole(gc, "admin")
		}
	}
	return false
}

func getPageSize(size *int) int {
	if size == nil || *size <= 0 {
		return 20
	}
	if *size > 100 {
		return 100
	}
	return *size
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getSortOrder(order *SortOrder) string {
	if order == nil || *order == SortOrderDesc {
		return "desc"
	}
	return "asc"
}

func convertPropertyInput(input *CreatePropertyInput) entity.Property {
	prop := entity.Property{
		ID:          uuid.New(),
		Name:        input.Name,
		DisplayName: input.DisplayName,
		DataType:    entity.DataType(input.DataType),
		IsRequired:  getBool(input.IsRequired),
		IsArray:     getBool(input.IsArray),
		IsUnique:    getBool(input.IsUnique),
		Metadata:    input.Metadata,
	}

	if input.DefaultValue != nil {
		prop.DefaultValue = input.DefaultValue
	}

	if input.ValidationRules != nil {
		prop.ValidationRules = convertValidationRulesInput(input.ValidationRules)
	}

	if input.Order != nil {
		prop.Order = *input.Order
	}

	return prop
}

func convertUpdatePropertyInput(input *UpdatePropertyInput) entity.Property {
	prop := entity.Property{}

	if input.ID != nil {
		prop.ID = *input.ID
	} else {
		prop.ID = uuid.New()
	}

	if input.Name != nil {
		prop.Name = *input.Name
	}
	if input.DisplayName != nil {
		prop.DisplayName = *input.DisplayName
	}
	if input.DataType != nil {
		prop.DataType = entity.DataType(*input.DataType)
	}
	if input.IsRequired != nil {
		prop.IsRequired = *input.IsRequired
	}
	if input.IsArray != nil {
		prop.IsArray = *input.IsArray
	}
	if input.IsUnique != nil {
		prop.IsUnique = *input.IsUnique
	}
	if input.DefaultValue != nil {
		prop.DefaultValue = input.DefaultValue
	}
	if input.ValidationRules != nil {
		prop.ValidationRules = convertValidationRulesInput(input.ValidationRules)
	}
	if input.Metadata != nil {
		prop.Metadata = input.Metadata
	}
	if input.Order != nil {
		prop.Order = *input.Order
	}

	return prop
}

func convertValidationRulesInput(input *ValidationRulesInput) *entity.ValidationRules {
	if input == nil {
		return nil
	}

	rules := &entity.ValidationRules{}

	if input.Pattern != nil {
		rules.Pattern = input.Pattern
	}
	if input.MinLength != nil {
		rules.MinLength = input.MinLength
	}
	if input.MaxLength != nil {
		rules.MaxLength = input.MaxLength
	}
	if input.MinValue != nil {
		rules.MinValue = input.MinValue
	}
	if input.MaxValue != nil {
		rules.MaxValue = input.MaxValue
	}
	if input.EnumValues != nil {
		rules.EnumValues = input.EnumValues
	}
	if input.CustomRules != nil {
		rules.CustomRules = input.CustomRules
	}

	return rules
}

func convertLinkConstraintsInput(input *LinkConstraintsInput) entity.LinkConstraints {
	constraints := entity.LinkConstraints{
		IsRequired:      getBool(input.IsRequired),
		CascadeDelete:   getBool(input.CascadeDelete),
		PreventDelete:   getBool(input.PreventDelete),
		UniquePerSource: getBool(input.UniquePerSource),
		UniquePerTarget: getBool(input.UniquePerTarget),
	}

	if input.ValidationRules != nil {
		constraints.ValidationRules = input.ValidationRules
	}

	return constraints
}

func getBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// encodeCursor encodes a cursor for pagination
func encodeCursor(timestamp time.Time, id uuid.UUID) string {
	data := fmt.Sprintf("%d:%s", timestamp.Unix(), id.String())
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Error definitions
var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)