package api

import (
	"database/sql"
	"os"

	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
)

func Login2FAHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body Login2FARequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.TempToken == "" || body.Code == "" {
			return c.Status(400).JSON(fiber.Map{"error": "missing token or code"})
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return c.Status(500).JSON(fiber.Map{"error": "JWT misconfigured"})
		}

		claims, err := security.ParseJWT(body.TempToken, secret)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid temp token"})
		}

		typ, _ := claims["typ"].(string)
		if typ != "2fa" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token type"})
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token payload"})
		}
		userID := int64(userIDFloat)

		var (
			username     string
			email        string
			totpSecret   string
			twoFAEnabled bool
		)

		err = db.QueryRow(`
			SELECT username, email, totp_secret, twofa_enabled
			FROM users
			WHERE id = $1
		`, userID).Scan(&username, &email, &totpSecret, &twoFAEnabled)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		if !twoFAEnabled || totpSecret == "" {
			return c.Status(400).JSON(fiber.Map{"error": "2fa not enabled"})
		}

		// Validate TOTP code (e.g. using pquerna/otp/totp under security.ValidateTOTP)
		if !security.ValidateTOTP(body.Code, totpSecret) {
			return c.Status(400).JSON(fiber.Map{"error": "invalid 2fa code"})
		}

		accessToken, err := createAccessToken(userID, email, username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "token error"})
		}

		return c.JSON(fiber.Map{
			"token": accessToken,
			"user":  buildUserResponse(userID, username, email, true, true),
		})
	}
}
