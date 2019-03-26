package cep

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	if isValid("29.315-732") {
		t.Errorf("Expected no valid CEP")
	}
	if !isValid("29315732") {
		t.Errorf("Expected a valid CEP")
	}
}

func TestSanitize(t *testing.T) {
	got := sanitize("29.315-732iX")
	expected := "29315732"
	if got != expected {
		t.Errorf("Expected '%s', got '%s'", got, expected)
	}
}
