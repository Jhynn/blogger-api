package security

import "golang.org/x/crypto/bcrypt"

// Hash returns a hash of the given text, otherwise an error.
func Hash(text string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
}

// CompareHashes returns an error if the raw string (the non-hashed one),
// which will be hashed, is different from the given hash; returns nil otherwise.
func CompareHashes(raw, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
}
