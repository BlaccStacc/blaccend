package api

import (
	"os"
	"time"

	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// helpers (foloseste camelCase)
func buildUserResponse(id int64, username, email string, emailVerified, twoFAEnabled bool) UserResponse {
	return UserResponse{
		ID:            id,
		Username:      username,
		Email:         email,
		EmailVerified: emailVerified,
		TwoFAEnabled:  twoFAEnabled,
	}
}

func createAccessToken(id int64, email, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fiber.NewError(fiber.StatusInternalServerError, "JWT missing")
	}

	claims := jwt.MapClaims{
		"user_id":  id,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"typ":      "access",
	}

	return security.SignJWT(claims, secret)
}

func createTemp2FAToken(id int64, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fiber.NewError(fiber.StatusInternalServerError, "JWT missing")
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(5 * time.Minute).Unix(), // expira in 5 min
		"typ":     "2fa",
	}

	return security.SignJWT(claims, secret)
}
