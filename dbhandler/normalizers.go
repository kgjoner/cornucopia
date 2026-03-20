package dbhandler

import (
	"regexp"
	"strings"
)

// normalizeMapKeys recursively normalizes all keys in a map, including nested maps and slices
func normalizeMapKeys(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		normalizedMap := make(map[string]interface{})
		for key, value := range v {
			normalizedKey := strings.ReplaceAll(key, "_", "")
			normalizedMap[normalizedKey] = normalizeMapKeys(value)
		}
		return normalizedMap

	case []interface{}:
		normalizedSlice := make([]interface{}, len(v))
		for i, item := range v {
			normalizedSlice[i] = normalizeMapKeys(item)
		}
		return normalizedSlice

	default:
		return v
	}
}

// normalizeTimestamps fixes timestamp format issues in the data
func normalizeTimestamps(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		normalizedMap := make(map[string]interface{})
		for key, value := range v {
			if str, ok := value.(string); ok {
				normalizedMap[key] = normalizeTimeString(str)
			} else {
				normalizedMap[key] = normalizeTimestamps(value)
			}
		}
		return normalizedMap

	case []interface{}:
		normalizedSlice := make([]interface{}, len(v))
		for i, item := range v {
			normalizedSlice[i] = normalizeTimestamps(item)
		}
		return normalizedSlice

	default:
		return v
	}
}

// ISO 8601 timestamp regex pattern
var timestampRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?$`)

// normalizeTimeString fixes individual time strings
func normalizeTimeString(str string) string {
	// Use regex to match ISO 8601 timestamp without timezone
	// Pattern: YYYY-MM-DDTHH:MM:SS[.ffffff]
	if timestampRegex.MatchString(str) {
		// Add Z to indicate UTC timezone
		return str + "Z"
	}
	return str
}