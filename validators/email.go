package validators

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string, required bool) bool {
	if !required {
		return true
	}

	return emailRegex.MatchString(email)
}
