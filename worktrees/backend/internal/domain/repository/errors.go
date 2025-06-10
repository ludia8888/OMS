package repository

import "errors"

// Common repository errors
var (
	// ErrCacheMiss indicates that the requested item was not found in cache
	ErrCacheMiss = errors.New("cache miss")
	
	// ErrNotFound indicates that the requested item was not found
	ErrNotFound = errors.New("not found")
	
	// ErrAlreadyExists indicates that an item with the same unique constraint already exists
	ErrAlreadyExists = errors.New("already exists")
	
	// ErrInvalidInput indicates that the provided input is invalid
	ErrInvalidInput = errors.New("invalid input")
	
	// ErrOptimisticLock indicates that the item was modified by another process
	ErrOptimisticLock = errors.New("optimistic lock failure")
)