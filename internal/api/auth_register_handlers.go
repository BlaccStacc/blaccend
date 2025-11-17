package api

import (
	"database/sql"
	"log"
	"time"

	//"time"
	"github.com/BlaccStacc/blaccend/internal/mail"
	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
)

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
			INSERT INTO users (username, email, password_hash, email_verified, verify_token, verify_expires_at, date_registered)
			VALUES ($1, $2, $3, FALSE, $4, $5, $6)
			RETURNING id
		`, body.Username, body.Email, passwordHash, verifyToken, verifyExpires, time.Now().UTC()).Scan(&id)
		if err != nil {
			// TODO: you can check for unique constraint on email and return 409
			return c.Status(400).JSON(fiber.Map{"error": "could not create user"})
		}

		if err := mail.SendVerificationEmail(body.Email, verifyToken); err != nil {
			log.Printf("SendVerificationEmail failed: %v", err)
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
