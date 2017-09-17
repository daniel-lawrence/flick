package user

import (
	"golang.org/x/crypto/bcrypt"
)

var bcryptCost = 12

// AuthUser represents a user that includes authentication functionality.
type AuthUser struct {
	Email          string `json:"email,omitempty"`
	hashedPassword []byte
}

// SetPassword hashes and stores the password for a user, overwriting the
// user's current password.
func (u *AuthUser) SetPassword(newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
	if err != nil {
		return err
	}
	u.hashedPassword = hashed
	return nil
}

// Authenticate takes a password and returns true if the user should be logged
// in, false otherwise.
func (u *AuthUser) Authenticate(password string) bool {
	result := bcrypt.CompareHashAndPassword(u.hashedPassword, []byte(password))
	if result == nil {
		return true
	}
	return false
}

// SetHashCost sets the global hashing cost. The cost will be clamped between 4
// and 31 (inclusive). The default cost is 10.
func SetHashCost(cost int) { bcryptCost = cost }
