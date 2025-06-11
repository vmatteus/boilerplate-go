package validator_test

import (
	"testing"

	"github.com/your-org/boilerplate-go/pkg/validator"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.com", true},
		{"", false},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"test@.com", false},
	}

	for _, test := range tests {
		result := validator.IsValidEmail(test.email)
		if result != test.valid {
			t.Errorf("IsValidEmail(%q) = %v, want %v", test.email, result, test.valid)
		}
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"John Doe", true},
		{"Mary-Jane", true},
		{"O'Connor", true},
		{"José", true},
		{"", false},
		{"   ", false},
		{"123", false},
		{"John@Doe", false},
	}

	for _, test := range tests {
		result := validator.IsValidName(test.name)
		if result != test.valid {
			t.Errorf("IsValidName(%q) = %v, want %v", test.name, result, test.valid)
		}
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"  hello world  ", 20, "hello world"},
		{"very long string", 5, "very "},
		{"", 10, ""},
		{"   ", 10, ""},
	}

	for _, test := range tests {
		result := validator.SanitizeString(test.input, test.maxLength)
		if result != test.expected {
			t.Errorf("SanitizeString(%q, %d) = %q, want %q", test.input, test.maxLength, result, test.expected)
		}
	}
}
