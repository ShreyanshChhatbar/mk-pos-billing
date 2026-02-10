package utils

import "strings"

// Comma-separated query parameter utilities
//
// Example 1: Simple string slice
//   tags := utils.CommaSeparatedStrings{}
//   tags.UnmarshalParam("tag1,tag2,tag3") // []string{"tag1", "tag2", "tag3"}
//
// Example 2: Enum/constant types (recommended)
//   statuses := utils.ParseCommaSeparatedEnum[constants.TaskStatus](c.Query("status"))
//   // Input: "PENDING,IN_PROGRESS" -> []constants.TaskStatus{PENDING, IN_PROGRESS}
//
// Example 3: In handler
//   filters := service.Filters{
//       Statuses: utils.ParseCommaSeparatedEnum[constants.TaskStatus](c.Query("status")),
//   }

// CommaSeparatedStrings is a custom type to handle comma-separated string values in query parameters
// Example: ?status=PENDING,IN_PROGRESS,COMPLETED
type CommaSeparatedStrings []string

// UnmarshalParam implements the binding interface for Gin to parse comma-separated values
func (c *CommaSeparatedStrings) UnmarshalParam(param string) error {
	if param == "" {
		*c = []string{}
		return nil
	}

	parts := strings.Split(param, ",")
	values := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}

	*c = values
	return nil
}

// CommaSeparatedEnum is a generic type for comma-separated enum values with validation
// T is the enum type (e.g., constants.TaskStatus)
type CommaSeparatedEnum[T ~string] []T

// UnmarshalParam parses comma-separated values and converts them to the enum type
func (c *CommaSeparatedEnum[T]) UnmarshalParam(param string) error {
	if param == "" {
		*c = []T{}
		return nil
	}

	parts := strings.Split(param, ",")
	values := make([]T, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, T(trimmed))
		}
	}

	*c = values
	return nil
}

// ParseCommaSeparatedEnum is a helper to parse comma-separated enum query parameters
// Usage: statuses := utils.ParseCommaSeparatedEnum[constants.TaskStatus](c.Query("status"))
func ParseCommaSeparatedEnum[T ~string](param string) []T {
	if param == "" {
		return []T{}
	}

	parts := strings.Split(param, ",")
	values := make([]T, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, T(trimmed))
		}
	}

	return values
}
