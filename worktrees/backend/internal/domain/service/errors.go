package service

import "errors"

// Service configuration errors
var (
	ErrMissingRepository           = errors.New("repository is required")
	ErrMissingObjectTypeRepository = errors.New("object type repository is required")
	ErrMissingCache               = errors.New("cache is required")
	ErrMissingEventPublisher      = errors.New("event publisher is required")
	ErrMissingLogger              = errors.New("logger is required")
)

// Business logic errors
var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrResourceNotFound    = errors.New("resource not found")
	ErrResourceExists      = errors.New("resource already exists")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrValidationFailed   = errors.New("validation failed")
	ErrConcurrentUpdate   = errors.New("concurrent update detected")
)