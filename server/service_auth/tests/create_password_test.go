package tests

import ("testing"
	sharedCrypto "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/crypto"	
)

/**
 * Test create password hash with salt
 */
func TestCreatePasswordHashWithSalt(t *testing.T) {
	// Input
	salt := "salt_test"
	password := []string{
		"admin.acme@example.com", 
		"alice.acme@example.com", 
		"bob.acme@example.com",
		"admin.beta@example.com",
		"charlie.beta@example.com",
	}
	// Expected output
	passwordHashWithSalt := []string{}
	for _, pwd := range password {
		passwordHash := sharedCrypto.HashPasswordWithSalt(pwd, salt)
		passwordHashWithSalt = append(passwordHashWithSalt, passwordHash)
	}
	// Log output
	for i, pwd := range password {
		t.Logf("Password: %s, Salt: %s, \n\tHash: %s\n", pwd, salt, passwordHashWithSalt[i])
	}
}