package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// aici se adauga rutele basically
func RegisterRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/health", GetHealth)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "hello",
		})
	})

	app.Get("/users/:id", GetUser(db))

	//app.Get("/auth")
}
