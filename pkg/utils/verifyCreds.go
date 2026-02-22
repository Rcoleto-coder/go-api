package utils

import "regexp"

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	passwordRegex := regexp.MustCompile(`^.{6,}$`) // at least 6 characters
	return passwordRegex.MatchString(password)
}