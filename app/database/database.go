package database

import (
	"ayo-baca-buku/app/config"
	"ayo-baca-buku/app/database/seeders"
	"ayo-baca-buku/app/models"
	"ayo-baca-buku/app/util/logger"
	"log"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
)

func NewDatabase(zLogger *zap.Logger) (*gorm.DB, error) {
	appConfig, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := appConfig.DB_SOURCE
	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	dbDebug := appConfig.DB_DEBUG
	gormLevel := gormLogger.Info

	if dbDebug != true {
		gormLevel = gormLogger.Silent
	}
	dbLogger := &logger.GormLogger{
		ZapLogger: zLogger,
		LogLevel:  gormLogger.Info,
	}
	db.Logger = dbLogger.LogMode(gormLevel)
	return db, nil
}

func RunMigration(DB *gorm.DB) {
	logger := logger.GetLogger()
	err := DB.AutoMigrate(
		&models.User{},
		&models.UserBook{},
		&models.ReadingActivity{},
	)

	if err != nil {
		logger.Fatal("Failed to migrate...")
	}

	logger.Info("Migrated Successfully")
}

func RunSeeder(DB *gorm.DB) {
	logger := logger.GetLogger()
	seeders.SeedUser(DB)
	logger.Info("Seeder Successfully")
}
