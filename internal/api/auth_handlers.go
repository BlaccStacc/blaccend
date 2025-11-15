package api

import (
	"database/sql"
	"os"
	"time"

	//"time"
	"github.com/BlaccStacc/blaccend/internal/mail"
	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

type Login2FARequest struct {
	TempToken string `json:"temp_token"`
	Code      string `json:"code"`
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

// GET /auth/verify-email?token=...
func VerifyEmailHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return c.Status(400).JSON(fiber.Map{"error": "missing token"})
		}

		res, err := db.Exec(`
			UPDATE users
			SET email_verified = TRUE,
			    verify_token = NULL,
			    verify_expires_at = NULL
			WHERE verify_token = $1
			  AND verify_expires_at > NOW()
		`, token)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		rows, err := res.RowsAffected()
		if err != nil || rows == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		return c.JSON(fiber.Map{"message": "email verified"})
	}
}

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

// POST /auth/login/2fa
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
