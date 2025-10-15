package password

import "golang.org/x/crypto/bcrypt"

// Hash returns a bcrypt hash of the password.
func Hash(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// Compare checks if the plain password matches the hashed password.
func Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
