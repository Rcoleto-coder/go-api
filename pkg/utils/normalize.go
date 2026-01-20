package utils

import "strings"

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NormalizePassword(password string) string {
	return strings.TrimSpace(password)
}
