// backend/internal/proj/c2c/service/marketplace_helpers.go
package service

// abs возвращает абсолютное значение числа
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// strPtr возвращает указатель на строку
func strPtr(s string) *string {
	return &s
}

// contains проверяет наличие строки в массиве
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
