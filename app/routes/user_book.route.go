package routes

import (
	"ayo-baca-buku/app/controllers" // Import the actual controllers package
	// "ayo-baca-buku/app/middlewares" // Placeholder for auth middleware if needed

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserBookRoutes(app *fiber.App, DB *gorm.DB) {
	// Instantiate the actual UserBookController
	userBookController := controllers.NewUserBookController(DB)

	// Group routes for /userbooks
	// Apply middleware here if needed, e.g., middlewares.AuthJWTMiddleware()
	userBookRoutes := app.Group("/userbooks")

	userBookRoutes.Post("/", userBookController.CreateUserBook)
	userBookRoutes.Get("/", userBookController.GetAllUserBooks)
	userBookRoutes.Get("/:id", userBookController.GetUserBookByID)
	userBookRoutes.Put("/:id", userBookController.UpdateUserBook)
	userBookRoutes.Delete("/:id", userBookController.DeleteUserBook) // Soft delete
}
