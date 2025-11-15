package security

// generates tokens and stuff

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func NewRandomToken(nBytes int) (string, error) {
	b, err := NewRandomBytes(nBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewURLSafeToken(nBytes int) (string, error) {
	b, err := NewRandomBytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func NewNumericCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}

	b, err := NewRandomBytes(length)
	if err != nil {
		return "", err
	}

	digits := make([]byte, length)
	for i := 0; i < length; i++ {
		digits[i] = '0' + (b[i] % 10)
	}

	return string(digits), nil
}
