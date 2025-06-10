package service

import (
	"github.com/openfoundry/oms/internal/domain/repository"
	"github.com/openfoundry/oms/internal/infrastructure/cache"
	"github.com/openfoundry/oms/internal/infrastructure/messaging"
	"go.uber.org/zap"
)

// ObjectTypeServiceConfig holds configuration for ObjectTypeService
type ObjectTypeServiceConfig struct {
	Repository     repository.ObjectTypeRepository
	Cache          cache.CacheService
	EventPublisher messaging.EventPublisher
	Logger         *zap.Logger
}

// Validate validates the configuration
func (c ObjectTypeServiceConfig) Validate() error {
	if c.Repository == nil {
		return ErrMissingRepository
	}
	if c.Cache == nil {
		return ErrMissingCache
	}
	if c.EventPublisher == nil {
		return ErrMissingEventPublisher
	}
	if c.Logger == nil {
		return ErrMissingLogger
	}
	return nil
}

// LinkTypeServiceConfig holds configuration for LinkTypeService
type LinkTypeServiceConfig struct {
	Repository         repository.LinkTypeRepository
	ObjectTypeRepo     repository.ObjectTypeRepository
	Cache              cache.CacheService
	EventPublisher     messaging.EventPublisher
	Logger             *zap.Logger
}

// Validate validates the configuration
func (c LinkTypeServiceConfig) Validate() error {
	if c.Repository == nil {
		return ErrMissingRepository
	}
	if c.ObjectTypeRepo == nil {
		return ErrMissingObjectTypeRepository
	}
	if c.Cache == nil {
		return ErrMissingCache
	}
	if c.EventPublisher == nil {
		return ErrMissingEventPublisher
	}
	if c.Logger == nil {
		return ErrMissingLogger
	}
	return nil
}