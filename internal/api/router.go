package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Get("/health", GetHealth)
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "hello"})
	})

	app.Get("/users/:id", GetUser(db))

	// AUTH
	app.Post("/auth/register", RegisterHandler(db))
	app.Post("/auth/login", LoginHandler(db))        // step 1 login
	app.Post("/auth/login/2fa", Login2FAHandler(db)) // step 2 login with TOTP
	app.Get("/auth/verify-email", VerifyEmailHandler(db))

	// AUTHENTICATED ROUTES
	protected := app.Group("", AuthMiddleware()) // require JWT
	protected.Get("/auth/me", MeHandler())

	// 🔹 TOTP setup (authenticator app)
	protected.Get("/auth/2fa/setup", TwoFASetupHandler(db))
	protected.Post("/auth/2fa/confirm", TwoFAConfirmHandler(db))

	// 🔹 GARAGE STORAGE ROUTES (nou)
	RegisterGarageRoutes(protected, db)
}
