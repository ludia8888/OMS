package entity

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

// Property represents a property of an object type
type Property struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"displayName"`
	DataType     DataType               `json:"dataType"`
	Required     bool                   `json:"required"`
	Unique       bool                   `json:"unique"`
	Indexed      bool                   `json:"indexed"`
	DefaultValue interface{}            `json:"defaultValue,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Validators   []Validator            `json:"validators,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// DataType represents the data type of a property
type DataType string

const (
	DataTypeString    DataType = "STRING"
	DataTypeNumber    DataType = "NUMBER"
	DataTypeBoolean   DataType = "BOOLEAN"
	DataTypeDate      DataType = "DATE"
	DataTypeDateTime  DataType = "DATETIME"
	DataTypeArray     DataType = "ARRAY"
	DataTypeObject    DataType = "OBJECT"
	DataTypeReference DataType = "REFERENCE"
)

// IsValid checks if the data type is valid
func (dt DataType) IsValid() bool {
	switch dt {
	case DataTypeString, DataTypeNumber, DataTypeBoolean,
		DataTypeDate, DataTypeDateTime, DataTypeArray,
		DataTypeObject, DataTypeReference:
		return true
	default:
		return false
	}
}

// Validator represents a validation rule for a property
type Validator struct {
	Type  ValidatorType `json:"type"`
	Value interface{}   `json:"value"`
}

// ValidatorType represents the type of validator
type ValidatorType string

const (
	ValidatorMinLength ValidatorType = "minLength"
	ValidatorMaxLength ValidatorType = "maxLength"
	ValidatorPattern   ValidatorType = "pattern"
	ValidatorMin       ValidatorType = "min"
	ValidatorMax       ValidatorType = "max"
	ValidatorEnum      ValidatorType = "enum"
	ValidatorFormat    ValidatorType = "format"
)

// IsValid checks if the validator type is valid
func (vt ValidatorType) IsValid() bool {
	switch vt {
	case ValidatorMinLength, ValidatorMaxLength, ValidatorPattern,
		ValidatorMin, ValidatorMax, ValidatorEnum, ValidatorFormat:
		return true
	default:
		return false
	}
}

// Validate validates the property definition
func (p *Property) Validate() error {
	if p.Name == "" {
		return ErrInvalidName
	}

	if !isValidPropertyName(p.Name) {
		return ErrInvalidPropertyNameFormat
	}

	if p.DisplayName == "" {
		return ErrRequiredField("displayName")
	}

	if !p.DataType.IsValid() {
		return ErrInvalidDataType(string(p.DataType))
	}

	// Validate validators
	for _, v := range p.Validators {
		if err := p.validateValidator(v); err != nil {
			return err
		}
	}

	// Validate default value if provided
	if p.DefaultValue != nil {
		if err := p.validateDefaultValue(); err != nil {
			return err
		}
	}

	return nil
}

// validateValidator validates a single validator
func (p *Property) validateValidator(v Validator) error {
	if !v.Type.IsValid() {
		return fmt.Errorf("invalid validator type: %s", v.Type)
	}

	// Validate validator value based on type
	switch v.Type {
	case ValidatorMinLength, ValidatorMaxLength:
		if p.DataType != DataTypeString {
			return fmt.Errorf("%s validator only applies to string type", v.Type)
		}
		if _, ok := v.Value.(float64); !ok {
			return fmt.Errorf("invalid value for %s validator", v.Type)
		}

	case ValidatorMin, ValidatorMax:
		if p.DataType != DataTypeNumber {
			return fmt.Errorf("%s validator only applies to number type", v.Type)
		}
		if _, ok := v.Value.(float64); !ok {
			return fmt.Errorf("invalid value for %s validator", v.Type)
		}

	case ValidatorPattern:
		if p.DataType != DataTypeString {
			return fmt.Errorf("pattern validator only applies to string type")
		}
		pattern, ok := v.Value.(string)
		if !ok {
			return fmt.Errorf("invalid pattern value")
		}
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}

	case ValidatorEnum:
		// Enum can apply to various types
		if _, ok := v.Value.([]interface{}); !ok {
			return fmt.Errorf("enum validator value must be an array")
		}
	}

	return nil
}

// validateDefaultValue validates the default value against the property type
func (p *Property) validateDefaultValue() error {
	switch p.DataType {
	case DataTypeString:
		if _, ok := p.DefaultValue.(string); !ok {
			return fmt.Errorf("default value must be a string for string type")
		}

	case DataTypeNumber:
		switch p.DefaultValue.(type) {
		case float64, int, int64, float32:
			// Valid number types
		default:
			return fmt.Errorf("default value must be a number for number type")
		}

	case DataTypeBoolean:
		if _, ok := p.DefaultValue.(bool); !ok {
			return fmt.Errorf("default value must be a boolean for boolean type")
		}

	case DataTypeDate, DataTypeDateTime:
		if _, ok := p.DefaultValue.(string); !ok {
			return fmt.Errorf("default value must be a string for date/datetime type")
		}

	case DataTypeArray:
		if _, ok := p.DefaultValue.([]interface{}); !ok {
			return fmt.Errorf("default value must be an array for array type")
		}

	case DataTypeObject:
		if _, ok := p.DefaultValue.(map[string]interface{}); !ok {
			return fmt.Errorf("default value must be an object for object type")
		}
	}

	return nil
}

// ValidateValue validates a value against the property definition
func (p *Property) ValidateValue(value interface{}) error {
	// Check required
	if p.Required && value == nil {
		return fmt.Errorf("required property %s is missing", p.Name)
	}

	// If not required and value is nil, it's valid
	if value == nil {
		return nil
	}

	// Validate type
	if err := p.validateValueType(value); err != nil {
		return err
	}

	// Apply validators
	for _, validator := range p.Validators {
		if err := applyValidator(validator, value, p.DataType); err != nil {
			return fmt.Errorf("validation failed for %s: %w", p.Name, err)
		}
	}

	return nil
}

// validateValueType validates that the value matches the expected type
func (p *Property) validateValueType(value interface{}) error {
	switch p.DataType {
	case DataTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("value must be a string for property %s", p.Name)
		}

	case DataTypeNumber:
		switch value.(type) {
		case float64, int, int64, float32:
			// Valid number types
		default:
			return fmt.Errorf("value must be a number for property %s", p.Name)
		}

	case DataTypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("value must be a boolean for property %s", p.Name)
		}

	case DataTypeArray:
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("value must be an array for property %s", p.Name)
		}

	case DataTypeObject:
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("value must be an object for property %s", p.Name)
		}

	case DataTypeReference:
		// Reference can be a string (ID) or an object
		switch value.(type) {
		case string, map[string]interface{}:
			// Valid reference types
		default:
			return fmt.Errorf("value must be a string or object for reference property %s", p.Name)
		}
	}

	return nil
}

// isValidPropertyName checks if the property name is valid
func isValidPropertyName(name string) bool {
	// Property names must start with lowercase letter and contain only alphanumeric and underscore
	pattern := `^[a-z][a-zA-Z0-9_]*$`
	matched, _ := regexp.MatchString(pattern, name)
	return matched && len(name) <= 64
}