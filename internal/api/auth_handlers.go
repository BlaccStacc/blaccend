package api

import (
	"database/sql"
	"os"
	"time"

	//"time"
	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// structuri
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	TwoFAEnabled  bool   `json:"twofa_enabled"`
}

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
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // expira dupa 24h
		"typ":      "access",
	}

	return security.SignJWT(claims, secret)
}

func createTemp2FAToken(id int62, email string) (string, error) {
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

// handlers (foloseste PascalCase)

// POST /auth/register
func RegisterHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body RegisterRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.Email == "" || body.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "email and password required"})
		}

		passwordHash, err := security.HashPassword(body.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "hashing failed"})
		}

		verifyToken, err := security.NewRandomToken(32)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "token error"})
		}
		verifyExpires := time.Now().Add(24 * time.Hour)

		var id int64
		err = db.QueryRow(`
			INSERT INTO users (username, email, password_hash, email_verified, verify_token, verify_expires_at)
			VALUES ($1, $2, $3, FALSE, $4, $5)
			RETURNING id
		`, body.Username, body.Email, passwordHash, verifyToken, verifyExpires).Scan(&id)
		if err != nil {
			// TODO: you can check for unique constraint on email and return 409
			return c.Status(400).JSON(fiber.Map{"error": "could not create user"})
		}

		if err := mail.SendVerificationEmail(body.Email, verifyToken); err != nil {
			// you might still keep the user and just log this
			// but for now return 500
			return c.Status(500).JSON(fiber.Map{"error": "failed to send verification email"})
		}

		return c.Status(201).JSON(buildUserResponse(id, body.Username, body.Email, false, false))

	}
}
