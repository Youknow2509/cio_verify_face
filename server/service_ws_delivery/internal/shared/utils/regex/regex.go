package regex

import "regexp"

// Check is a valid email address
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}