package routes

import (
	"ayo-baca-buku/app/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, DB *gorm.DB) {
	userController := controllers.NewUserController(DB)

	// Group routes for users
	userRoutes := app.Group("/users")

	userRoutes.Get("/", userController.GetAllUsers)
	userRoutes.Post("/", userController.CreateUser) // Added CreateUser route
	userRoutes.Get("/:id", userController.GetUserById) // Added GetUserById route
	userRoutes.Put("/:id", userController.UpdateUser) // Added UpdateUser route
	userRoutes.Delete("/:id", userController.DeleteUser) // Added DeleteUser (hard delete) route
	userRoutes.Patch("/:id/soft-delete", userController.SoftDeleteUser) // Added SoftDeleteUser route (using PATCH for partial update semantics)

	// Example of a protected route group if you have middleware for auth
	// protectedUserRoutes := app.Group("/users", middlewares.AuthMiddleware) // Assuming you have an AuthMiddleware
	// protectedUserRoutes.Get("/", userController.GetAllUsers)
	// ... other protected routes
}
