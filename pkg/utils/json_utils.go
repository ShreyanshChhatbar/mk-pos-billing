package utils

import (
	"encoding/json"
	"fmt"
)

func DecodeJSONRawAny(raw json.RawMessage) (interface{}, bool) {
	if len(raw) == 0 {
		return nil, false
	}
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, false
	}
	return v, true
}

func ParseJSONRawString(raw json.RawMessage) string {
	val, ok := DecodeJSONRawAny(raw)
	if !ok {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ParseJSONRawStringPtr(raw json.RawMessage) *string {
	val := ParseJSONRawString(raw)
	if val == "" {
		return nil
	}
	return &val
}

func ParseJSONRawFloat(raw json.RawMessage) float64 {
	val, ok := DecodeJSONRawAny(raw)
	if !ok {
		return 0
	}
	return ParseFloat(val)
}

func ParseJSONRawUint(raw json.RawMessage) uint64 {
	val, ok := DecodeJSONRawAny(raw)
	if !ok {
		return 0
	}
	return ParseUint(val)
}
