package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openfoundry/oms/internal/domain/entity"
	"github.com/openfoundry/oms/internal/domain/repository"
	"github.com/openfoundry/oms/internal/infrastructure/cache"
	"github.com/openfoundry/oms/internal/infrastructure/messaging"
	"go.uber.org/zap"
)

// ObjectTypeService handles business logic for object types
type ObjectTypeService struct {
	repo      repository.ObjectTypeRepository
	cache     cache.CacheService
	publisher messaging.EventPublisher
	logger    *zap.Logger
}

// NewObjectTypeService creates a new object type service
func NewObjectTypeService(
	repo repository.ObjectTypeRepository,
	cache cache.CacheService,
	publisher messaging.EventPublisher,
	logger *zap.Logger,
) *ObjectTypeService {
	return &ObjectTypeService{
		repo:      repo,
		cache:     cache,
		publisher: publisher,
		logger:    logger,
	}
}

// CreateObjectTypeInput represents input for creating an object type
type CreateObjectTypeInput struct {
	Name         string                         `json:"name"`
	DisplayName  string                         `json:"displayName"`
	Description  *string                        `json:"description"`
	Category     *string                        `json:"category"`
	Tags         []string                       `json:"tags"`
	Properties   []PropertyInput                `json:"properties"`
	Metadata     map[string]interface{}         `json:"metadata"`
}

// PropertyInput represents input for creating a property
type PropertyInput struct {
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"displayName"`
	DataType     entity.DataType        `json:"dataType"`
	Required     bool                   `json:"required"`
	Unique       bool                   `json:"unique"`
	Indexed      bool                   `json:"indexed"`
	DefaultValue interface{}            `json:"defaultValue,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Validators   []entity.Validator     `json:"validators,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// CreateObjectType creates a new object type
func (s *ObjectTypeService) CreateObjectType(ctx context.Context, input CreateObjectTypeInput, userID string) (*entity.ObjectType, error) {
	s.logger.Info("Creating object type", zap.String("name", input.Name), zap.String("user", userID))

	// Check if name already exists
	existing, _ := s.repo.GetByName(ctx, input.Name)
	if existing != nil {
		return nil, entity.ErrObjectTypeNameExists
	}

	// Build properties
	properties := make([]entity.Property, len(input.Properties))
	for i, propInput := range input.Properties {
		properties[i] = entity.Property{
			ID:           uuid.New(),
			Name:         propInput.Name,
			DisplayName:  propInput.DisplayName,
			DataType:     propInput.DataType,
			Required:     propInput.Required,
			Unique:       propInput.Unique,
			Indexed:      propInput.Indexed,
			DefaultValue: propInput.DefaultValue,
			Description:  propInput.Description,
			Validators:   propInput.Validators,
			Metadata:     propInput.Metadata,
		}
	}

	// Create object type entity
	objectType := &entity.ObjectType{
		ID:          uuid.New(),
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Description: input.Description,
		Category:    input.Category,
		Tags:        input.Tags,
		Properties:  properties,
		Metadata:    input.Metadata,
		Version:     1,
		IsDeleted:   false,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
		UpdatedAt:   time.Now(),
		UpdatedBy:   userID,
	}

	// Validate object type
	if err := objectType.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Save to repository
	if err := s.repo.Create(ctx, objectType); err != nil {
		s.logger.Error("Failed to create object type", zap.Error(err))
		return nil, fmt.Errorf("failed to create object type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache(ctx, objectType.ID)

	// Publish event
	event := messaging.Event{
		ID:        uuid.New().String(),
		Type:      messaging.EventObjectTypeCreated,
		EntityID:  objectType.ID.String(),
		Actor:     userID,
		Timestamp: time.Now(),
		Data:      objectType,
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	s.logger.Info("Object type created successfully", zap.String("id", objectType.ID.String()))
	return objectType, nil
}

// GetByID retrieves an object type by ID
func (s *ObjectTypeService) GetByID(ctx context.Context, id uuid.UUID) (*entity.ObjectType, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("object_type:%s", id.String())
	var cached *entity.ObjectType
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil && cached != nil {
		return cached, nil
	}

	// Get from repository
	objectType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, objectType, 5*time.Minute)

	return objectType, nil
}

// GetByName retrieves an object type by name
func (s *ObjectTypeService) GetByName(ctx context.Context, name string) (*entity.ObjectType, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("object_type:name:%s", name)
	var cached *entity.ObjectType
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil && cached != nil {
		return cached, nil
	}

	// Get from repository
	objectType, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, objectType, 5*time.Minute)

	return objectType, nil
}

// UpdateObjectTypeInput represents input for updating an object type
type UpdateObjectTypeInput struct {
	DisplayName *string                        `json:"displayName,omitempty"`
	Description *string                        `json:"description,omitempty"`
	Category    *string                        `json:"category,omitempty"`
	Tags        []string                       `json:"tags,omitempty"`
	Properties  []PropertyInput                `json:"properties,omitempty"`
	Metadata    map[string]interface{}         `json:"metadata,omitempty"`
}

// UpdateObjectType updates an existing object type
func (s *ObjectTypeService) UpdateObjectType(ctx context.Context, id uuid.UUID, input UpdateObjectTypeInput, userID string) (*entity.ObjectType, error) {
	s.logger.Info("Updating object type", zap.String("id", id.String()), zap.String("user", userID))

	// Get existing object type
	objectType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.DisplayName != nil {
		objectType.DisplayName = *input.DisplayName
	}
	if input.Description != nil {
		objectType.Description = input.Description
	}
	if input.Category != nil {
		objectType.Category = input.Category
	}
	if input.Tags != nil {
		objectType.Tags = input.Tags
	}
	if input.Properties != nil {
		// Convert property inputs
		properties := make([]entity.Property, len(input.Properties))
		for i, propInput := range input.Properties {
			properties[i] = entity.Property{
				ID:           uuid.New(),
				Name:         propInput.Name,
				DisplayName:  propInput.DisplayName,
				DataType:     propInput.DataType,
				Required:     propInput.Required,
				Unique:       propInput.Unique,
				Indexed:      propInput.Indexed,
				DefaultValue: propInput.DefaultValue,
				Description:  propInput.Description,
				Validators:   propInput.Validators,
				Metadata:     propInput.Metadata,
			}
		}
		objectType.Properties = properties
	}
	if input.Metadata != nil {
		objectType.Metadata = input.Metadata
	}

	// Update metadata
	objectType.IncrementVersion()
	objectType.SetUpdatedBy(userID)

	// Validate
	if err := objectType.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Save to repository
	if err := s.repo.Update(ctx, objectType); err != nil {
		s.logger.Error("Failed to update object type", zap.Error(err))
		return nil, fmt.Errorf("failed to update object type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache(ctx, objectType.ID)

	// Publish event
	event := messaging.Event{
		ID:        uuid.New().String(),
		Type:      messaging.EventObjectTypeUpdated,
		EntityID:  objectType.ID.String(),
		Actor:     userID,
		Timestamp: time.Now(),
		Data:      objectType,
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	s.logger.Info("Object type updated successfully", zap.String("id", objectType.ID.String()))
	return objectType, nil
}

// DeleteObjectType soft deletes an object type
func (s *ObjectTypeService) DeleteObjectType(ctx context.Context, id uuid.UUID, userID string) error {
	s.logger.Info("Deleting object type", zap.String("id", id.String()), zap.String("user", userID))

	// Check if object type exists
	objectType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check for dependencies (e.g., instances, link types)

	// Soft delete
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete object type", zap.Error(err))
		return fmt.Errorf("failed to delete object type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache(ctx, id)

	// Publish event
	event := messaging.Event{
		ID:        uuid.New().String(),
		Type:      messaging.EventObjectTypeDeleted,
		EntityID:  id.String(),
		Actor:     userID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"objectTypeId": id.String(),
			"name":        objectType.Name,
		},
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	s.logger.Info("Object type deleted successfully", zap.String("id", id.String()))
	return nil
}

// List retrieves a list of object types based on filter
func (s *ObjectTypeService) List(ctx context.Context, filter repository.ObjectTypeFilter) ([]*entity.ObjectType, error) {
	return s.repo.List(ctx, filter)
}

// Search searches for object types
func (s *ObjectTypeService) Search(ctx context.Context, query string, limit int) ([]*entity.ObjectType, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("object_types:search:%s:%d", query, limit)
	var cached []*entity.ObjectType
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil && cached != nil {
		return cached, nil
	}

	// Search in repository
	results, err := s.repo.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// Cache the results
	_ = s.cache.Set(ctx, cacheKey, results, 2*time.Minute)

	return results, nil
}

// CompareVersions compares two versions of an object type
func (s *ObjectTypeService) CompareVersions(ctx context.Context, id uuid.UUID, v1, v2 int) (*repository.VersionDiff, error) {
	return s.repo.CompareVersions(ctx, id, v1, v2)
}

// Count counts object types based on filter
func (s *ObjectTypeService) Count(ctx context.Context, filter repository.ObjectTypeFilter) (int64, error) {
	return s.repo.Count(ctx, filter)
}

// GetVersion retrieves a specific version of an object type
func (s *ObjectTypeService) GetVersion(ctx context.Context, id uuid.UUID, version int) (*repository.ObjectTypeVersion, error) {
	return s.repo.GetVersion(ctx, id, version)
}

// ListVersions lists all versions of an object type
func (s *ObjectTypeService) ListVersions(ctx context.Context, id uuid.UUID) ([]*repository.ObjectTypeVersion, error) {
	return s.repo.ListVersions(ctx, id)
}

// invalidateCache invalidates cache entries for an object type
func (s *ObjectTypeService) invalidateCache(ctx context.Context, id uuid.UUID) {
	_ = s.cache.Delete(ctx, fmt.Sprintf("object_type:%s", id.String()))
	_ = s.cache.InvalidatePattern(ctx, "object_types:*")
}