package routes

import (
	"bosshire.com/handlers"
	"bosshire.com/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	//Non public API
	api := app.Group("/api", middleware.Protected())

	api.Get("/jobs", handlers.ViewJobs)
	api.Post("/jobs", handlers.PostJob)
	api.Post("/jobs/:id/apply", handlers.ApplyJob)
	api.Get("/applications", handlers.ReviewApplications)
	api.Put("/applications/:id", handlers.ProcessApplication)
}
