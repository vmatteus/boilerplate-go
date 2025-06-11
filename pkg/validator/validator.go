package validator

import (
	"regexp"
	"strings"
)

// IsValidEmail validates if the email format is correct
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidName validates if the name is not empty and contains only allowed characters
func IsValidName(name string) bool {
	name = strings.TrimSpace(name)
	if len(name) == 0 || len(name) > 255 {
		return false
	}

	// Allow letters, spaces, hyphens, and apostrophes
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	return nameRegex.MatchString(name)
}

// SanitizeString removes leading/trailing whitespace and limits length
func SanitizeString(s string, maxLength int) string {
	s = strings.TrimSpace(s)
	if len(s) > maxLength {
		return s[:maxLength]
	}
	return s
}
