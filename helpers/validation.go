package helpers

import (
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	// Trim spaces
	email = strings.TrimSpace(email)

	// Regular expression for email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
