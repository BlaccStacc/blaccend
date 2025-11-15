package api

import (
	"database/sql"

	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// POST /auth/login
func LoginHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body LoginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		var (
			id            int64
			username      string
			email         string
			passwordHash  string
			emailVerified bool
			twoFAEnabled  bool
		)

		err := db.QueryRow(`
			SELECT id, username, email, password_hash, email_verified, twofa_enabled
			FROM users
			WHERE email = $1
		`, body.Email).Scan(&id, &username, &email, &passwordHash, &emailVerified, &twoFAEnabled)

		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		if !emailVerified {
			return c.Status(403).JSON(fiber.Map{"error": "email not verified"})
		}

		if !security.VerifyPassword(passwordHash, body.Password) {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		if twoFAEnabled {
			tempToken, err := createTemp2FAToken(id, email)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "token error"})
			}

			return c.JSON(fiber.Map{
				"requires_2fa": true,
				"temp_token":   tempToken,
				"user":         buildUserResponse(id, username, email, emailVerified, twoFAEnabled),
			})
		}

		accesToken, err := createAccessToken(id, email, username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "token error"})
		}

		return c.JSON(fiber.Map{
			"token": accesToken,
			"user":  buildUserResponse(id, username, email, emailVerified, twoFAEnabled),
		})
	}
}

// GET /auth/me
func MeHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		idFloat, _ := claims["user_id"].(float64)
		username, _ := claims["username"].(string)
		email, _ := claims["email"].(string)

		resp := UserResponse{
			ID:       int64(idFloat),
			Username: username,
			Email:    email,
		}

		return c.JSON(resp)
	}
}
