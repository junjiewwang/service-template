package internal

import (
	"fmt"
)

// ValidateNotEmpty validates that a string is not empty
func ValidateNotEmpty(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidatePositive validates that a number is positive
func ValidatePositive(value int, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got %d", fieldName, value)
	}
	return nil
}

// ValidateInRange validates that a number is within a range
func ValidateInRange(value, min, max int, fieldName string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %d and %d, got %d", fieldName, min, max, value)
	}
	return nil
}

// ValidateOneOf validates that a value is one of the allowed values
func ValidateOneOf(value string, allowed []string, fieldName string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %v, got %s", fieldName, allowed, value)
}
