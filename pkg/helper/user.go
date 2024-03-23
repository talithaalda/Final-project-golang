package helper

import (
	"math"
	"regexp"
	"time"
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
func CountAge(dob time.Time) int {
	age := time.Since(dob).Hours()
	age = age / 24 / 365
	age = math.Round(age)
	return int(age)
}