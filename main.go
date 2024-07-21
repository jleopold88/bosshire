package main

import (
	db "bosshire.com/db"
	"bosshire.com/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	db.Connect()
	// simple cors
	app.Use(cors.New())
	routes.Setup(app)
	app.Listen(":3000")
}
