package main

import (
    "log"
	
    "github.com/gofiber/fiber/v2"
    "github.com/BlaccStacc/blaccend/internal/api"
    "github.com/BlaccStacc/blaccend/internal/config"
    "github.com/BlaccStacc/blaccend/internal/db"
)

func main() {
	cfg := config.Load()

	// postgres connection?
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db conect failed: %v", err)
	}
	defer database.Close() // defer schedules a func to run after the surrounding func returns, no matter how
	//basically defer runs last


	// create fiber app ce pula mea e fiber- web framework pt go gen un fel de node.js lolz
	app := fiber.New()

	api.RegisterRoutes(app, database)

	log.Printf("server running pe port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}