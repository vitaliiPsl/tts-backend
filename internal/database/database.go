package database

import (
	"os"
	"time"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	dsn := os.Getenv("DATABASE_CONNECTION_STRING")

	logger.Logger.Info("Connecting to the database....", "dsn", dsn)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Error("Error connecting to database", "Error message", err)
		panic(nil)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Logger.Error("Error getting DB from GORM", "Error message", err)
		panic(nil)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logger.Logger.Info("Connected to the database.")

	logger.Logger.Info("Migrating models...")
	DB.AutoMigrate(&user.User{})
	logger.Logger.Info("Migrated models.")
}
