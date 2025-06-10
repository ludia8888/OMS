package entity

import (
	"time"

	"github.com/google/uuid"
)

// LinkType represents a relationship between two object types
type LinkType struct {
	ID                 uuid.UUID              `json:"id"`
	Name               string                 `json:"name"`
	DisplayName        string                 `json:"displayName"`
	SourceObjectTypeID uuid.UUID              `json:"sourceObjectTypeId"`
	TargetObjectTypeID uuid.UUID              `json:"targetObjectTypeId"`
	Cardinality        Cardinality            `json:"cardinality"`
	Description        *string                `json:"description,omitempty"`
	Properties         []Property             `json:"properties,omitempty"`
	Metadata           map[string]interface{} `json:"metadata"`
	Version            int                    `json:"version"`
	IsDeleted          bool                   `json:"-"`
	CreatedAt          time.Time              `json:"createdAt"`
	CreatedBy          string                 `json:"createdBy"`
	UpdatedAt          time.Time              `json:"updatedAt"`
	UpdatedBy          string                 `json:"updatedBy"`
}

// Cardinality represents the cardinality of a relationship
type Cardinality string

const (
	CardinalityOneToOne   Cardinality = "ONE_TO_ONE"
	CardinalityOneToMany  Cardinality = "ONE_TO_MANY"
	CardinalityManyToMany Cardinality = "MANY_TO_MANY"
)

// IsValid checks if the cardinality is valid
func (c Cardinality) IsValid() bool {
	switch c {
	case CardinalityOneToOne, CardinalityOneToMany, CardinalityManyToMany:
		return true
	default:
		return false
	}
}

// Validate validates the link type
func (lt *LinkType) Validate() error {
	if lt.Name == "" {
		return ErrInvalidName
	}

	if !isValidName(lt.Name) {
		return ErrInvalidNameFormat
	}

	if lt.DisplayName == "" {
		return ErrRequiredField("displayName")
	}

	if lt.SourceObjectTypeID == uuid.Nil {
		return ErrRequiredField("sourceObjectTypeId")
	}

	if lt.TargetObjectTypeID == uuid.Nil {
		return ErrRequiredField("targetObjectTypeId")
	}

	if !lt.Cardinality.IsValid() {
		return ErrInvalidCardinality(string(lt.Cardinality))
	}

	// Validate properties if any
	propertyNames := make(map[string]bool)
	for _, prop := range lt.Properties {
		if propertyNames[prop.Name] {
			return ErrDuplicateProperty(prop.Name)
		}
		propertyNames[prop.Name] = true

		if err := prop.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// IncrementVersion increments the version number
func (lt *LinkType) IncrementVersion() {
	lt.Version++
	lt.UpdatedAt = time.Now()
}

// SetUpdatedBy sets the updated by field and timestamp
func (lt *LinkType) SetUpdatedBy(userID string) {
	lt.UpdatedBy = userID
	lt.UpdatedAt = time.Now()
}

// IsSelfReferencing checks if the link type is self-referencing
func (lt *LinkType) IsSelfReferencing() bool {
	return lt.SourceObjectTypeID == lt.TargetObjectTypeID
}

// GetInverseCardinality returns the inverse cardinality
func (lt *LinkType) GetInverseCardinality() Cardinality {
	switch lt.Cardinality {
	case CardinalityOneToOne:
		return CardinalityOneToOne
	case CardinalityOneToMany:
		return CardinalityManyToMany // Inverse of 1:N is N:1 (represented as M:N with constraint)
	case CardinalityManyToMany:
		return CardinalityManyToMany
	default:
		return lt.Cardinality
	}
}