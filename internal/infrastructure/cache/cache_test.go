package cache

import (
	"encoding/json"
	"testing"
)

func TestPHPValueToJSONCompatibleConvertsSequentialMapToSlice(t *testing.T) {
	input := map[interface{}]interface{}{
		int64(0): map[interface{}]interface{}{
			"category_id": int64(7),
			"discount":    float64(5),
		},
		int64(1): map[interface{}]interface{}{
			"category_id": int64(8),
			"discount":    float64(10),
		},
	}

	got := phpValueToJSONCompatible(input)

	list, ok := got.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", got)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}

	first, ok := list[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected first item to be map[string]interface{}, got %T", list[0])
	}
	if first["category_id"] != int64(7) {
		t.Fatalf("expected first category_id to be 7, got %#v", first["category_id"])
	}
}

func TestPHPValueToJSONCompatibleConvertsNestedAssociativeMap(t *testing.T) {
	input := map[interface{}]interface{}{
		"id": int64(23),
		"product_categories_discounts": map[interface{}]interface{}{
			int64(0): map[interface{}]interface{}{
				"category_id": int64(10),
				"value":       "15",
			},
		},
	}

	got := phpValueToJSONCompatible(input)

	obj, ok := got.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T", got)
	}

	discounts, ok := obj["product_categories_discounts"].([]interface{})
	if !ok {
		t.Fatalf("expected product_categories_discounts to be []interface{}, got %T", obj["product_categories_discounts"])
	}
	if len(discounts) != 1 {
		t.Fatalf("expected 1 discount item, got %d", len(discounts))
	}

	if _, ok := discounts[0].(map[string]interface{}); !ok {
		t.Fatalf("expected nested discount item to be map[string]interface{}, got %T", discounts[0])
	}

	if _, err := json.Marshal(got); err != nil {
		t.Fatalf("expected converted value to marshal to JSON, got error: %v", err)
	}
}
