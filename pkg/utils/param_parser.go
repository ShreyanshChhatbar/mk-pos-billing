package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseUint64Param extracts and parses a uint64 parameter from the Gin context
// Returns the parsed value and an error if parsing fails
func ParseUint64Param(c *gin.Context, paramName string) (uint64, error) {
	paramStr := c.Param(paramName)
	if paramStr == "" {
		return 0, fmt.Errorf("parameter '%s' is required", paramName)
	}

	value, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be a valid number", paramName)
	}

	return value, nil
}

// ParseUint64 parses a string to uint64
// Returns the parsed value and an error if parsing fails
func ParseUint64(s string) (uint64, error) {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	return value, nil
}
