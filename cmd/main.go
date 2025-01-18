package main

import (
	"ayo-baca-buku/app/database"
	"ayo-baca-buku/app/routes"
	"ayo-baca-buku/app/util/logger"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	gormLogger "gorm.io/gorm/logger"
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
	zLogger := logger.NewLogger()
	defer zLogger.Sync()

	DB, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	dbLogger := &logger.GormLogger{
		ZapLogger: zLogger,
		LogLevel:  gormLogger.Info,
	}
	DB.Logger = dbLogger.LogMode(gormLogger.Info)

	database.RunMigration(DB)
	database.RunSeeder(DB)

	app := fiber.New()
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zLogger,
	}))
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

	go func() {
		// Memberikan sedikit jeda untuk memastikan server sudah berjalan
		time.Sleep(100 * time.Millisecond)
		banner := `
    _______ __             
   / ____(_) /_  ___  _____
  / /_  / / __ \/ _ \/ ___/
 / __/ / / /_/ /  __/ /    
/_/   /_/_.___/\___/_/`
		fmt.Println(banner)
		fmt.Println("\nPress Ctrl+C to shutdown server")
	}()

	if err := app.Listen(":3000"); err != nil {
		zLogger.Error("Server failed to start", zap.Error(err))
		os.Exit(1)
	}
}
