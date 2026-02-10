package constants

import "strings"

// StringEnum is a constraint for string-based enum types
type StringEnum interface {
	~string
}

// IsValidEnum checks if a value matches any of the allowed enum values
func IsValidEnum[T StringEnum](value string, allowedValues ...T) bool {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	for _, allowed := range allowedValues {
		if normalized == string(allowed) {
			return true
		}
	}
	return false
}

// AllValues is a generic function that returns all values for an enum type
// Usage: constants.AllValues(FrequencyTypeOnce, FrequencyTypeDaily, ...)
func AllValues[T StringEnum](values ...T) []T {
	return values
}
