package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "ayo-baca-buku/docs"
)

func InitSwagger(app *fiber.App) {
	app.Get("/docs/*", swagger.HandlerDefault)
}
