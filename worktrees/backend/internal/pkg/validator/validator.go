package validator

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var (
	// Object type name pattern: must start with letter, contain only alphanumeric and underscore
	objectTypeNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	
	// Property name pattern: must start with lowercase letter, contain only alphanumeric and underscore
	propertyNamePattern = regexp.MustCompile(`^[a-z][a-zA-Z0-9_]*$`)
	
	// Email pattern
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	// URL pattern
	urlPattern = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
)

// ValidateObjectTypeName validates an object type name
func ValidateObjectTypeName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	
	if len(name) > 64 {
		return fmt.Errorf("name must not exceed 64 characters")
	}
	
	if !objectTypeNamePattern.MatchString(name) {
		return fmt.Errorf("name must start with a letter and contain only alphanumeric characters and underscores")
	}
	
	// Check for reserved words
	reserved := []string{"system", "meta", "internal", "private", "public"}
	lowerName := strings.ToLower(name)
	for _, r := range reserved {
		if lowerName == r {
			return fmt.Errorf("name '%s' is reserved", name)
		}
	}
	
	return nil
}

// ValidatePropertyName validates a property name
func ValidatePropertyName(name string) error {
	if name == "" {
		return fmt.Errorf("property name cannot be empty")
	}
	
	if len(name) > 64 {
		return fmt.Errorf("property name must not exceed 64 characters")
	}
	
	if !propertyNamePattern.MatchString(name) {
		return fmt.Errorf("property name must start with a lowercase letter and contain only alphanumeric characters and underscores")
	}
	
	// Check for reserved property names
	reserved := []string{"id", "createdAt", "updatedAt", "createdBy", "updatedBy", "version"}
	for _, r := range reserved {
		if name == r {
			return fmt.Errorf("property name '%s' is reserved", name)
		}
	}
	
	return nil
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	
	if !emailPattern.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// ValidateURL validates a URL
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	
	if !urlPattern.MatchString(url) {
		return fmt.Errorf("invalid URL format")
	}
	
	return nil
}

// SanitizeString sanitizes a string to prevent XSS attacks
func SanitizeString(input string) string {
	// HTML escape the string
	sanitized := html.EscapeString(input)
	
	// Remove any null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	
	return sanitized
}

// SanitizeTags sanitizes a list of tags
func SanitizeTags(tags []string) []string {
	sanitized := make([]string, 0, len(tags))
	seen := make(map[string]bool)
	
	for _, tag := range tags {
		// Sanitize each tag
		cleanTag := SanitizeString(tag)
		
		// Skip empty tags
		if cleanTag == "" {
			continue
		}
		
		// Skip duplicates
		if seen[cleanTag] {
			continue
		}
		
		seen[cleanTag] = true
		sanitized = append(sanitized, cleanTag)
	}
	
	return sanitized
}

// ValidatePageSize validates pagination page size
func ValidatePageSize(size int) (int, error) {
	if size <= 0 {
		return 20, nil // Default
	}
	
	if size > 100 {
		return 0, fmt.Errorf("page size cannot exceed 100")
	}
	
	return size, nil
}

// ValidateSortOrder validates sort order
func ValidateSortOrder(order string) (string, error) {
	order = strings.ToLower(order)
	
	switch order {
	case "asc", "desc", "":
		return order, nil
	default:
		return "", fmt.Errorf("sort order must be 'asc' or 'desc'")
	}
}

// ValidateSortBy validates sort field
func ValidateSortBy(field string, allowedFields []string) (string, error) {
	if field == "" {
		return "", nil
	}
	
	for _, allowed := range allowedFields {
		if field == allowed {
			return field, nil
		}
	}
	
	return "", fmt.Errorf("invalid sort field: %s", field)
}