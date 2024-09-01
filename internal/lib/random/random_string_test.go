// random_string_test.go
package random

import (
	"testing"
)

// TestNewRandomString_Length проверяет, что длина сгенерированной строки соответствует ожидаемой.
func TestNewRandomString_Length(t *testing.T) {
	for length := 0; length < 11; length++ {
		result := NewRandomString(length)
		if len(result) != length {
			t.Errorf("Expected string length of %d, but got %d", length, len(result))
		}
	}
}

// TestNewRandomString_Charset проверяет, что сгенерированная строка содержит только разрешенные символы.
func TestNewRandomString_Charset(t *testing.T) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := NewRandomString(100) // Generate a longer string for better charset coverage

	for _, char := range result {
		if !contains(charset, char) {
			t.Errorf("Generated string contains invalid character: %c", char)
		}
	}
}

// TestNewRandomString_Uniqueness проверяет, что функция генерирует уникальные строки при многократном вызове.
func TestNewRandomString_Uniqueness(t *testing.T) {
	length := 10
	iterations := 1000
	results := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		result := NewRandomString(length)
		if results[result] {
			t.Errorf("Duplicate string generated: %s", result)
		}
		results[result] = true
	}
}

// contains проверяет, содержится ли символ в строке.
func contains(s string, c rune) bool {
	for _, char := range s {
		if char == c {
			return true
		}
	}
	return false
}
