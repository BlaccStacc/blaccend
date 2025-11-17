package api

import (
	"database/sql"
	"os"

	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// POST /auth/login/2fa
// Body: { "temp_token": "...", "code": "123456" }
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

// GET /auth/2fa/setup   (protected; requires AuthMiddleware)
// Returns a new TOTP secret + otpauth URL for QR code.
func TwoFASetupHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		idFloat, _ := claims["user_id"].(float64)
		email, _ := claims["email"].(string)
		if idFloat == 0 || email == "" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token payload"})
		}
		userID := int64(idFloat)

		// Generate a new TOTP secret for this user
		secret, err := security.GenerateTOTPSecret()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to generate 2fa secret"})
		}

		// Build otpauth URL (used by authenticator apps / QR code)
		otpauthURL, err := security.BuildOtpauthURL(secret, "Blaccend", email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to build otpauth url"})
		}

		// Store secret in DB but keep twofa_enabled = false until confirmed
		_, err = db.Exec(`
			UPDATE users
			SET totp_secret = $1,
			    twofa_enabled = FALSE
			WHERE id = $2
		`, secret, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		return c.JSON(fiber.Map{
			"secret":      secret,
			"otpauth_url": otpauthURL,
		})
	}
}

// POST /auth/2fa/confirm   (protected; requires AuthMiddleware)
// Body: { "code": "123456" }
// Validates the code against the stored secret and flips twofa_enabled = true.
func TwoFAConfirmHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		idFloat, _ := claims["user_id"].(float64)
		if idFloat == 0 {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token payload"})
		}
		userID := int64(idFloat)

		var body struct {
			Code string `json:"code"`
		}
		if err := c.BodyParser(&body); err != nil || body.Code == "" {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		var totpSecret string
		err := db.QueryRow(`
			SELECT totp_secret
			FROM users
			WHERE id = $1
		`, userID).Scan(&totpSecret)
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}
		if totpSecret == "" {
			return c.Status(400).JSON(fiber.Map{"error": "2fa not initialized"})
		}

		if !security.ValidateTOTP(body.Code, totpSecret) {
			return c.Status(400).JSON(fiber.Map{"error": "invalid 2fa code"})
		}

		_, err = db.Exec(`
			UPDATE users
			SET twofa_enabled = TRUE
			WHERE id = $1
		`, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		return c.JSON(fiber.Map{
			"message":       "2fa enabled",
			"twofa_enabled": true,
		})
	}
}
