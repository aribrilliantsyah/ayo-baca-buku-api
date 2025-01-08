package main

import (
	"ayo-baca-buku/app/config"
	"ayo-baca-buku/app/models"
	"ayo-baca-buku/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToPostgreSQL() (*gorm.DB, error) {
	appConfig, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := appConfig.DB_SOURCE
	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitializeMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.UserBook{},
		&models.ReadingActivity{},
	)
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	db, err := connectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	err = InitializeMigration(db)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	util.InitSwagger(app)

	app.Listen(": 3000")
}
