package entity

import (
	"time"

	"github.com/google/uuid"
)

// ObjectType represents a business object definition
type ObjectType struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"displayName"`
	Description  *string                `json:"description,omitempty"`
	Category     *string                `json:"category,omitempty"`
	Tags         []string               `json:"tags"`
	Properties   []Property             `json:"properties"`
	BaseDatasets []DatasetReference     `json:"baseDatasets,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
	Version      int                    `json:"version"`
	IsDeleted    bool                   `json:"-"`
	CreatedAt    time.Time              `json:"createdAt"`
	CreatedBy    string                 `json:"createdBy"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	UpdatedBy    string                 `json:"updatedBy"`
}

// DatasetReference represents a reference to a base dataset
type DatasetReference struct {
	DatasetRID string `json:"datasetRid"`
	Name       string `json:"name"`
}

// Validate validates the object type
func (ot *ObjectType) Validate() error {
	if ot.Name == "" {
		return ErrInvalidName
	}

	if !isValidName(ot.Name) {
		return ErrInvalidNameFormat
	}

	if ot.DisplayName == "" {
		return ErrRequiredField("displayName")
	}

	// Validate properties
	propertyNames := make(map[string]bool)
	for _, prop := range ot.Properties {
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
func (ot *ObjectType) IncrementVersion() {
	ot.Version++
	ot.UpdatedAt = time.Now()
}

// SetUpdatedBy sets the updated by field and timestamp
func (ot *ObjectType) SetUpdatedBy(userID string) {
	ot.UpdatedBy = userID
	ot.UpdatedAt = time.Now()
}

// AddProperty adds a new property to the object type
func (ot *ObjectType) AddProperty(prop Property) error {
	// Check for duplicate
	for _, existing := range ot.Properties {
		if existing.Name == prop.Name {
			return ErrDuplicateProperty(prop.Name)
		}
	}

	// Validate property
	if err := prop.Validate(); err != nil {
		return err
	}

	ot.Properties = append(ot.Properties, prop)
	return nil
}

// RemoveProperty removes a property by name
func (ot *ObjectType) RemoveProperty(propertyName string) error {
	for i, prop := range ot.Properties {
		if prop.Name == propertyName {
			// Remove the property
			ot.Properties = append(ot.Properties[:i], ot.Properties[i+1:]...)
			return nil
		}
	}
	return ErrPropertyNotFound(propertyName)
}

// UpdateProperty updates an existing property
func (ot *ObjectType) UpdateProperty(propertyName string, updatedProp Property) error {
	for i, prop := range ot.Properties {
		if prop.Name == propertyName {
			// Validate the updated property
			if err := updatedProp.Validate(); err != nil {
				return err
			}

			// Update the property
			ot.Properties[i] = updatedProp
			return nil
		}
	}
	return ErrPropertyNotFound(propertyName)
}

// GetProperty returns a property by name
func (ot *ObjectType) GetProperty(propertyName string) (*Property, error) {
	for _, prop := range ot.Properties {
		if prop.Name == propertyName {
			return &prop, nil
		}
	}
	return nil, ErrPropertyNotFound(propertyName)
}

// HasTag checks if the object type has a specific tag
func (ot *ObjectType) HasTag(tag string) bool {
	for _, t := range ot.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag if it doesn't already exist
func (ot *ObjectType) AddTag(tag string) {
	if !ot.HasTag(tag) {
		ot.Tags = append(ot.Tags, tag)
	}
}

// RemoveTag removes a tag
func (ot *ObjectType) RemoveTag(tag string) {
	for i, t := range ot.Tags {
		if t == tag {
			ot.Tags = append(ot.Tags[:i], ot.Tags[i+1:]...)
			return
		}
	}
}