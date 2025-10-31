package random

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// RandomString generates a random string of length n
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// RandomInt generates a random string of length 6
func GenerateSixDigitOtp() int {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := 100000 + rng.Intn(900000) // 100000 -> 999999

	return otp
}

// GenerateUUID generates a random UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// Random password user
func RandomPassword() string {
	return fmt.Sprintf("%s-%s-%s-%s", RandomString(6), RandomString(6), RandomString(6), RandomString(6))
}
