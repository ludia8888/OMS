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

// LinkTypeService handles business logic for link types
type LinkTypeService struct {
	repo            repository.LinkTypeRepository
	objectTypeRepo  repository.ObjectTypeRepository
	cache           cache.CacheService
	publisher       messaging.EventPublisher
	logger          *zap.Logger
}

// NewLinkTypeService creates a new link type service
func NewLinkTypeService(config LinkTypeServiceConfig) (*LinkTypeService, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &LinkTypeService{
		repo:           config.Repository,
		objectTypeRepo: config.ObjectTypeRepo,
		cache:          config.Cache,
		publisher:      config.EventPublisher,
		logger:         config.Logger,
	}, nil
}

// CreateLinkTypeInput represents input for creating a link type
type CreateLinkTypeInput struct {
	Name               string                 `json:"name"`
	DisplayName        string                 `json:"displayName"`
	InverseDisplayName *string                `json:"inverseDisplayName"`
	Description        *string                `json:"description"`
	SourceObjectTypeID uuid.UUID              `json:"sourceObjectTypeId"`
	TargetObjectTypeID uuid.UUID              `json:"targetObjectTypeId"`
	Cardinality        entity.Cardinality     `json:"cardinality"`
	Properties         []entity.Property      `json:"properties"`
	Constraints        entity.LinkConstraints `json:"constraints"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// UpdateLinkTypeInput represents input for updating a link type
type UpdateLinkTypeInput struct {
	DisplayName        *string                 `json:"displayName"`
	InverseDisplayName *string                 `json:"inverseDisplayName"`
	Description        *string                 `json:"description"`
	Cardinality        *entity.Cardinality     `json:"cardinality"`
	Properties         *[]entity.Property      `json:"properties"`
	Constraints        *entity.LinkConstraints `json:"constraints"`
	Metadata           map[string]interface{}  `json:"metadata"`
}

// CreateLinkType creates a new link type
func (s *LinkTypeService) CreateLinkType(ctx context.Context, input CreateLinkTypeInput, userID string) (*entity.LinkType, error) {
	s.logger.Info("Creating link type",
		zap.String("name", input.Name),
		zap.String("source", input.SourceObjectTypeID.String()),
		zap.String("target", input.TargetObjectTypeID.String()),
		zap.String("user", userID))

	// Verify source and target object types exist
	if _, err := s.objectTypeRepo.GetByID(ctx, input.SourceObjectTypeID); err != nil {
		if err == repository.ErrNotFound {
			return nil, entity.ErrObjectTypeNotFound
		}
		return nil, fmt.Errorf("failed to verify source object type: %w", err)
	}

	if _, err := s.objectTypeRepo.GetByID(ctx, input.TargetObjectTypeID); err != nil {
		if err == repository.ErrNotFound {
			return nil, entity.ErrObjectTypeNotFound
		}
		return nil, fmt.Errorf("failed to verify target object type: %w", err)
	}

	// Check for circular reference if needed
	if input.SourceObjectTypeID == input.TargetObjectTypeID {
		// Self-referencing is allowed, but check if it would create issues
		if hasCircular, err := s.repo.CheckCircularReference(ctx, input.SourceObjectTypeID, input.TargetObjectTypeID); err != nil {
			return nil, fmt.Errorf("failed to check circular reference: %w", err)
		} else if hasCircular {
			return nil, entity.ErrCircularReference
		}
	}

	// Check if link type name already exists
	if existing, err := s.repo.GetByName(ctx, input.Name); err == nil && existing != nil {
		return nil, entity.ErrLinkTypeNameExists
	}

	// Create link type entity
	linkType := &entity.LinkType{
		ID:                 uuid.New(),
		Name:               input.Name,
		DisplayName:        input.DisplayName,
		InverseDisplayName: input.InverseDisplayName,
		Description:        input.Description,
		SourceObjectTypeID: input.SourceObjectTypeID,
		TargetObjectTypeID: input.TargetObjectTypeID,
		Cardinality:        input.Cardinality,
		Properties:         input.Properties,
		Constraints:        input.Constraints,
		Metadata:           input.Metadata,
		Version:            1,
		IsDeleted:          false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          userID,
		UpdatedBy:          userID,
	}

	// Validate
	if err := linkType.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Save to repository
	if err := s.repo.Create(ctx, linkType); err != nil {
		s.logger.Error("Failed to create link type", zap.Error(err))
		return nil, fmt.Errorf("failed to create link type: %w", err)
	}

	// Publish event
	event := messaging.Event{
		ID:        uuid.New().String(),
		Type:      messaging.EventLinkTypeCreated,
		EntityID:  linkType.ID.String(),
		Actor:     userID,
		Timestamp: time.Now(),
		Data:      linkType,
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	s.logger.Info("Link type created successfully", zap.String("id", linkType.ID.String()))
	return linkType, nil
}

// GetByID retrieves a link type by ID
func (s *LinkTypeService) GetByID(ctx context.Context, id uuid.UUID) (*entity.LinkType, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("link_type:%s", id.String())
	var cached entity.LinkType
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// Get from repository
	linkType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, linkType, 5*time.Minute)

	return linkType, nil
}

// GetByName retrieves a link type by name
func (s *LinkTypeService) GetByName(ctx context.Context, name string) (*entity.LinkType, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("link_type:name:%s", name)
	var cached entity.LinkType
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// Get from repository
	linkType, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, linkType, 5*time.Minute)

	return linkType, nil
}

// UpdateLinkType updates an existing link type
func (s *LinkTypeService) UpdateLinkType(ctx context.Context, id uuid.UUID, input UpdateLinkTypeInput, userID string) (*entity.LinkType, error) {
	s.logger.Info("Updating link type", zap.String("id", id.String()), zap.String("user", userID))

	// Get existing link type
	linkType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.DisplayName != nil {
		linkType.DisplayName = *input.DisplayName
	}
	if input.InverseDisplayName != nil {
		linkType.InverseDisplayName = input.InverseDisplayName
	}
	if input.Description != nil {
		linkType.Description = input.Description
	}
	if input.Cardinality != nil {
		linkType.Cardinality = *input.Cardinality
	}
	if input.Properties != nil {
		linkType.Properties = *input.Properties
	}
	if input.Constraints != nil {
		linkType.Constraints = *input.Constraints
	}
	if input.Metadata != nil {
		linkType.Metadata = input.Metadata
	}

	// Update metadata
	linkType.IncrementVersion()
	linkType.SetUpdatedBy(userID)

	// Validate
	if err := linkType.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Save to repository
	if err := s.repo.Update(ctx, linkType); err != nil {
		s.logger.Error("Failed to update link type", zap.Error(err))
		return nil, fmt.Errorf("failed to update link type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache(ctx, linkType.ID)

	// Publish event
	event := messaging.Event{
		ID:        uuid.New().String(),
		Type:      messaging.EventLinkTypeUpdated,
		EntityID:  linkType.ID.String(),
		Actor:     userID,
		Timestamp: time.Now(),
		Data:      linkType,
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	s.logger.Info("Link type updated successfully", zap.String("id", linkType.ID.String()))
	return linkType, nil
}

// DeleteLinkType soft deletes a link type
func (s *LinkTypeService) DeleteLinkType(ctx context.Context, id uuid.UUID, userID string) error {
	s.logger.Info("Deleting link type", zap.String("id", id.String()), zap.String("user", userID))

	// Check if link type exists
	linkType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check for dependencies (e.g., link instances)

	// Soft delete
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete link type", zap.Error(err))
		return fmt.Errorf("failed to delete link type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache(ctx, id)

	// Publish event
	event := messaging.Event{
		ID:        uuid.New().String(),
		Type:      messaging.EventLinkTypeDeleted,
		EntityID:  id.String(),
		Actor:     userID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"linkTypeId": id.String(),
			"name":       linkType.Name,
		},
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish event", zap.Error(err))
	}

	s.logger.Info("Link type deleted successfully", zap.String("id", id.String()))
	return nil
}

// List retrieves a list of link types based on filter
func (s *LinkTypeService) List(ctx context.Context, filter repository.LinkTypeFilter) ([]*entity.LinkType, error) {
	return s.repo.List(ctx, filter)
}

// Count counts link types based on filter
func (s *LinkTypeService) Count(ctx context.Context, filter repository.LinkTypeFilter) (int64, error) {
	return s.repo.Count(ctx, filter)
}

// GetBySourceObjectType retrieves link types by source object type
func (s *LinkTypeService) GetBySourceObjectType(ctx context.Context, objectTypeID uuid.UUID) ([]*entity.LinkType, error) {
	return s.repo.GetBySourceObjectType(ctx, objectTypeID)
}

// GetByTargetObjectType retrieves link types by target object type
func (s *LinkTypeService) GetByTargetObjectType(ctx context.Context, objectTypeID uuid.UUID) ([]*entity.LinkType, error) {
	return s.repo.GetByTargetObjectType(ctx, objectTypeID)
}

// GetByObjectTypes retrieves link types between two object types
func (s *LinkTypeService) GetByObjectTypes(ctx context.Context, sourceID, targetID uuid.UUID) ([]*entity.LinkType, error) {
	return s.repo.GetByObjectTypes(ctx, sourceID, targetID)
}

// CheckCircularReference checks if creating a link would result in a circular reference
func (s *LinkTypeService) CheckCircularReference(ctx context.Context, sourceID, targetID uuid.UUID) (bool, error) {
	return s.repo.CheckCircularReference(ctx, sourceID, targetID)
}

// invalidateCache invalidates cache entries for a link type
func (s *LinkTypeService) invalidateCache(ctx context.Context, id uuid.UUID) {
	_ = s.cache.Delete(ctx, fmt.Sprintf("link_type:%s", id.String()))
	_ = s.cache.InvalidatePattern(ctx, "link_types:*")
}