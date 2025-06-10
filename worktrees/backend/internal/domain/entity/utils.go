package entity

import (
	"fmt"
	"regexp"
)

// isValidName checks if the name is valid for object types and link types
func isValidName(name string) bool {
	// Names must start with letter and contain only alphanumeric and underscore
	pattern := `^[a-zA-Z][a-zA-Z0-9_]*$`
	matched, _ := regexp.MatchString(pattern, name)
	return matched && len(name) <= 64
}

// applyValidator applies a validator to a value
func applyValidator(validator Validator, value interface{}, dataType DataType) error {
	switch validator.Type {
	case ValidatorMinLength:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("value is not a string")
		}
		minLen, ok := validator.Value.(float64)
		if !ok {
			return fmt.Errorf("invalid minLength value")
		}
		if len(str) < int(minLen) {
			return fmt.Errorf("string length %d is less than minimum %d", len(str), int(minLen))
		}

	case ValidatorMaxLength:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("value is not a string")
		}
		maxLen, ok := validator.Value.(float64)
		if !ok {
			return fmt.Errorf("invalid maxLength value")
		}
		if len(str) > int(maxLen) {
			return fmt.Errorf("string length %d exceeds maximum %d", len(str), int(maxLen))
		}

	case ValidatorPattern:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("value is not a string")
		}
		pattern, ok := validator.Value.(string)
		if !ok {
			return fmt.Errorf("invalid pattern value")
		}
		matched, err := regexp.MatchString(pattern, str)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
		if !matched {
			return fmt.Errorf("value does not match pattern %s", pattern)
		}

	case ValidatorMin:
		var num float64
		switch v := value.(type) {
		case float64:
			num = v
		case int:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		default:
			return fmt.Errorf("value is not a number")
		}
		min, ok := validator.Value.(float64)
		if !ok {
			return fmt.Errorf("invalid min value")
		}
		if num < min {
			return fmt.Errorf("value %f is less than minimum %f", num, min)
		}

	case ValidatorMax:
		var num float64
		switch v := value.(type) {
		case float64:
			num = v
		case int:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		default:
			return fmt.Errorf("value is not a number")
		}
		max, ok := validator.Value.(float64)
		if !ok {
			return fmt.Errorf("invalid max value")
		}
		if num > max {
			return fmt.Errorf("value %f exceeds maximum %f", num, max)
		}

	case ValidatorEnum:
		enumValues, ok := validator.Value.([]interface{})
		if !ok {
			return fmt.Errorf("invalid enum values")
		}
		found := false
		for _, enumVal := range enumValues {
			if value == enumVal {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value is not in enum")
		}
	}

	return nil
}