package unstructured

import "errors"

type Unstructured struct {
	Object map[string]any
}

func NestedInt64(obj map[string]any, fields ...string) (int64, bool, error) {
	current := obj
	for i, f := range fields {
		val, ok := current[f]
		if !ok {
			return 0, false, nil
		}
		if i == len(fields)-1 {
			switch v := val.(type) {
			case int:
				return int64(v), true, nil
			case int64:
				return v, true, nil
			case float64:
				return int64(v), true, nil
			default:
				return 0, false, errors.New("unexpected type")
			}
		}
		next, ok := val.(map[string]any)
		if !ok {
			return 0, false, errors.New("unexpected structure")
		}
		current = next
	}
	return 0, false, nil
}

func NestedString(obj map[string]any, fields ...string) (string, bool, error) {
	current := obj
	for i, f := range fields {
		val, ok := current[f]
		if !ok {
			return "", false, nil
		}
		if i == len(fields)-1 {
			if s, ok := val.(string); ok {
				return s, true, nil
			}
			return "", false, errors.New("unexpected type")
		}
		next, ok := val.(map[string]any)
		if !ok {
			return "", false, errors.New("unexpected structure")
		}
		current = next
	}
	return "", false, nil
}
