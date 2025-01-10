package routes

import (
	"ayo-baca-buku/app/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, DB *gorm.DB) {
	userController := controllers.NewUserController(DB)

	app.Get("/users", userController.GetAllUsers)
}
