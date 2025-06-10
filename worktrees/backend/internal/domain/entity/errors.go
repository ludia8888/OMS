package entity

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	// Object Type errors
	ErrObjectTypeNotFound   = errors.New("object type not found")
	ErrObjectTypeNameExists = errors.New("object type name already exists")
	ErrInvalidObjectType    = errors.New("invalid object type")
	
	// Property errors
	ErrPropertyNotFound          = errors.New("property not found")
	ErrInvalidPropertyNameFormat = errors.New("property name must start with lowercase letter and contain only alphanumeric and underscore")
	
	// Link Type errors
	ErrLinkTypeNotFound   = errors.New("link type not found")
	ErrLinkTypeNameExists = errors.New("link type name already exists")
	ErrCircularReference  = errors.New("circular reference detected")
	
	// General validation errors
	ErrInvalidName       = errors.New("name is required")
	ErrInvalidNameFormat = errors.New("name must start with letter and contain only alphanumeric and underscore")
	ErrRequiredFieldMissing = errors.New("required field is missing")
)

// ErrRequiredField returns an error for a missing required field
func ErrRequiredField(fieldName string) error {
	return fmt.Errorf("%w: %s", ErrRequiredFieldMissing, fieldName)
}

// ErrDuplicateProperty returns an error for duplicate property
func ErrDuplicateProperty(propertyName string) error {
	return fmt.Errorf("duplicate property name: %s", propertyName)
}

// ErrPropertyNotFoundWithName returns an error for property not found
func ErrPropertyNotFound(propertyName string) error {
	return fmt.Errorf("%w: %s", ErrPropertyNotFound, propertyName)
}

// ErrInvalidDataType returns an error for invalid data type
func ErrInvalidDataType(dataType string) error {
	return fmt.Errorf("invalid data type: %s", dataType)
}

// ErrInvalidCardinality returns an error for invalid cardinality
func ErrInvalidCardinality(cardinality string) error {
	return fmt.Errorf("invalid cardinality: %s", cardinality)
}