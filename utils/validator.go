package utils

import (
	"errors"
	"unicode"
)

var (
	ErrWeakPassword = errors.New("password must be atleast 8 characters and should contain uppercase, lowercase, special characters, number")
)

func isPasswordValid(password string) bool {
	var (
		hasMinLen = false
		hasUpper = false
		hasLower = false
		hasNumber = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}