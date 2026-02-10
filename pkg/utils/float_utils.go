package utils

import (
	"fmt"
	"math"
)

// RoundToTwoDecimals rounds a float64 to 2 decimal places
// This is essential for monetary and quantity calculations to avoid floating-point precision errors
func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// FormatToTwoDecimals formats a float64 to a string with exactly 2 decimal places
func FormatToTwoDecimals(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// Subtract calculates the difference between two values with 2-decimal precision
// Returns: a - b, rounded to 2 decimals
func Subtract(a, b float64) float64 {
	return RoundToTwoDecimals(a - b)
}

// MultiplyAbsolute calculates the product of absolute value and multiplier with 2-decimal precision
// Returns: |value| × multiplier, rounded to 2 decimals
func MultiplyAbsolute(value, multiplier float64) float64 {
	return RoundToTwoDecimals(math.Abs(value) * multiplier)
}

// Multiply calculates the product of two values with 2-decimal precision
// Returns: a × b, rounded to 2 decimals
func Multiply(a, b float64) float64 {
	return RoundToTwoDecimals(a * b)
}

// FormatAsInteger formats a float64 to a string as a whole number (no decimals)
// Used for quantities that should be stored as integers
func FormatAsInteger(value float64) string {
	return fmt.Sprintf("%.0f", value)
}

// Percentage calculates the percentage of reviewed out of total
// Returns: (reviewed / total) × 100, rounded to 2 decimals
// Returns 0 if total is 0 to avoid division by zero
func Percentage(total, reviewed int) float64 {
	if total == 0 {
		return 0
	}
	return RoundToTwoDecimals((float64(reviewed) * 100) / float64(total))
}
