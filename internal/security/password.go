package security

// argon2id hashing and verification

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/BlaccStacc/blaccend/internal/security"
	"golang.org/x/crypto/argon2"
)

// params for argon2
const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

func HashPassword(password string) (string, error) {
	salt, err := security.NewRandomBytes(16)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

func VerifyPassword(encodedHash, password string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	// parts[4] = salt, parts[5] = hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// Recalculate hash
	computed := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, uint32(len(expectedHash)))

	return constantTimeCompare(expectedHash, computed)
}

func constantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var diff byte
	for i := 0; i < len(a); i++ {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}
