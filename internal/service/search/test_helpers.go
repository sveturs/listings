package search

// Test helper functions shared across all test files

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}

func ptrString(v string) *string {
	return &v
}

func generateLongString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

func generateStringSlice(length int) []string {
	result := make([]string, length)
	for i := range result {
		result[i] = "value"
	}
	return result
}
