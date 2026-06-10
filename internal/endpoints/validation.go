package endpoints

import "regexp"

// NOTE: string before @ string after then "." and string. this is the email regex
var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// NOTE: Lower and upper cases from a - z, numbers from 0 to 9 and it must be between 3 and 32 characters. this is the username regex
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

func isValidEmail(email string) bool{
	return emailRegex.MatchString(email)
}

func isValidUsername(username string) bool{
	return usernameRegex.MatchString(username)
}

func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 128{
		return false
	}
	hasLower := false
	hasDigit := false

	for _, r := range password {
		switch{
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}
	return hasLower && hasDigit
}
