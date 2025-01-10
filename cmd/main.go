package main

import (
	"ayo-baca-buku/app/database"
	"ayo-baca-buku/app/routes"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Ayo Baca Buku - API
// @version 1.0
// @description Ini adalah API - Ayo Baca Buku
// @contact.name Ari Ganteng
// @contact.email ariardiansyah.study@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	DB, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	database.RunMigration(DB)

	app := fiber.New()
	app.Static("/docs", "docs")
	app.Get("/docs/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.json",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Ayo Baca Buku - API",
		})
	})

	app.Get("/scalar", func(c *fiber.Ctx) error {
		html := fmt.Sprintf(`<!doctype html>
		<html lang="en">
			<head>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<title>Swagger API Reference - Scalar</title>
				<link rel="icon" type="image/svg+xml" href="https://docs.scalar.com/favicon.svg">
				<link rel="icon" type="image/png" href="https://docs.scalar.com/favicon.png">
			</head>
			<body>
				<script id="api-reference" data-url="%s"></script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
		</html>`, "/docs/swagger.json")

		return c.Type("html").Send([]byte(html))
	})

	routes.SetupUserRoutes(app, DB)

	app.Listen(": 3000")
}
