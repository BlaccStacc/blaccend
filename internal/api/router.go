package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	// CORS so the React app (5173) can talk to backend (8080)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	app.Get("/health", GetHealth)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "hello",
		})
	})

	app.Get("/users/:id", GetUser(db))

	// AUTH: these must be POST, not GET
	app.Post("/auth/register", RegisterHandler(db))
	app.Post("/auth/login", LoginHandler(db))
	app.Post("/auth/login/2fa", Login2FAHandler(db))

	// /auth/me is a GET (used by frontend to fetch current user)
	app.Get("/auth/me", MeHandler())
}
