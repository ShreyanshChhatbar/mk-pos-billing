package utils

func ToMap(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if ok {
				m[key] = args[i+1]
			}
		}
	}
	return m
}
