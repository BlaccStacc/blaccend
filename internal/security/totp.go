package security

import (
	//"time"

	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret() (string, error) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Blaccend", //to change??
		AccountName: "",
	})
	if err != nil {
		return "", err
	}

	return secret.Secret(), nil
}

func ValidateTOTP(code string, secret string) bool {
	return totp.Validate(code, secret)
}

func BuildOtpauthURL(secret, issuer, email string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Secret:      []byte(secret),
		Issuer:      issuer,
		AccountName: email,
	})
	if err != nil {
		return "", err
	}
	return key.URL(), nil
}
