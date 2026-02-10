package utils

import (
	"fmt"
	"strconv"
)

// ParseInt converts various types to int64
func ParseInt(value interface{}) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	case []uint8: // Request body bytes
		i, _ := strconv.ParseInt(string(v), 10, 64)
		return i
	default:
		// Attempt string conversion for other types
		str := fmt.Sprintf("%v", v)
		i, _ := strconv.ParseInt(str, 10, 64)
		return i
	}
}

// ParseUint converts various types to uint64
func ParseUint(value interface{}) uint64 {
	switch v := value.(type) {
	case int:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int8:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int16:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int32:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int64:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case uint:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	case float32:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case float64:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case string:
		i, _ := strconv.ParseUint(v, 10, 64)
		return i
	case []uint8:
		i, _ := strconv.ParseUint(string(v), 10, 64)
		return i
	default:
		str := fmt.Sprintf("%v", v)
		i, _ := strconv.ParseUint(str, 10, 64)
		return i
	}
}

// ParseFloat converts various types to float64
func ParseFloat(value interface{}) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case []uint8:
		f, _ := strconv.ParseFloat(string(v), 64)
		return f
	default:
		str := fmt.Sprintf("%v", v)
		f, _ := strconv.ParseFloat(str, 64)
		return f
	}
}

// ToUint converts uint64 to uint
func ToUint(value uint64) uint {
	return uint(value)
}
