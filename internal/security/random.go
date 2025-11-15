package security

import (
	"crypto/rand"
	"fmt"
)

// should theoretically return SECURELY generated random bytes lolz
func NewRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("generate random bytes: %w", err)
	}
	return b, nil
}
