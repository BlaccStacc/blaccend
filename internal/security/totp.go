package security

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateTOTPSecret generates a random base32 secret suitable for TOTP.
// Example: "JBSWY3DPEHPK3PXP"
func GenerateTOTPSecret() (string, error) {
	// 20 bytes = 160 bits, standard for TOTP
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	secret := enc.EncodeToString(buf)
	secret = strings.TrimRight(secret, "=")
	return secret, nil
}

// BuildOtpauthURL builds an otpauth:// URL that authenticator apps understand.
// Example: otpauth://totp/Issuer:email@example.com?secret=ABC...&issuer=Issuer&digits=6&period=30
func BuildOtpauthURL(secret, issuer, email string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("empty secret")
	}
	if issuer == "" {
		issuer = "Blaccend"
	}
	if email == "" {
		email = "user"
	}

	label := fmt.Sprintf("%s:%s", issuer, email)

	return fmt.Sprintf(
		"otpauth://totp/%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		url.PathEscape(label),
		secret,
		url.QueryEscape(issuer),
	), nil
}

// ValidateTOTP checks if the provided TOTP code is valid for the given secret at the current time.
func ValidateTOTP(code, secret string) bool {
	code = strings.TrimSpace(code)
	if code == "" || secret == "" {
		return false
	}

	ok, err := totp.ValidateCustom(
		code,
		secret,
		time.Now(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      1, // allow +/-1 step
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	if err != nil {
		return false
	}
	return ok
}
