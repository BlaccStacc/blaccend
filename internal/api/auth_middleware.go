package api

import (
	"os"
	"strings"

	"github.com/BlaccStacc/blaccend/internal/security"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Expect "Authorization: Bearer <token>"
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing auth header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid auth header"})
		}

		token := parts[1]

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return c.Status(500).JSON(fiber.Map{"error": "JWT misconfigured"})
		}

		claims, err := security.ParseJWT(token, secret)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// Attach claims to context
		c.Locals("user", claims)

		return c.Next()
	}
}
