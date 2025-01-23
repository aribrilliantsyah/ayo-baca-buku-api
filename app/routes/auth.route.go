package routes

import (
	"ayo-baca-buku/app/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, DB *gorm.DB) {
	authController := controllers.NewAuthController(DB)

	app.Post("/login", authController.Login)
	app.Post("/register", authController.Register)
}
