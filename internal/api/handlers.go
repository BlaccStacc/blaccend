package api

import (
	"database/sql"
    "github.com/gofiber/fiber/v2"
)

// func ExampleHandler(c *fiber.Ctx) error {
// 	return c.JSON(fiber.Map{
// 		"message": "example route",
// 	})
// }

func GetUser(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var username, email string
		err := db.QueryRow("SELECT username, email FROM users WHERE id=$1", id).Scan(&username, &email)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(404).JSON(fiber.Map{"error": "user not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "db error"})
		}

		return c.JSON(fiber.Map{
			"id": id,
			"username": username,
			"email": email,
		})
	}
}