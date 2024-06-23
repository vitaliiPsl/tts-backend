package database

import (
	"os"
	"time"
	"vitaliiPsl/synthesizer/internal/history"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/token"
	"vitaliiPsl/synthesizer/internal/users"
	"vitaliiPsl/synthesizer/internal/model"

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
		panic(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Logger.Error("Error getting DB from GORM", "Error message", err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logger.Logger.Info("Connected to the database.")

	logger.Logger.Info("Migrating models...")
	DB.AutoMigrate(&users.User{}, &token.Token{}, &history.HistoryRecord{}, &model.Model{})
	logger.Logger.Info("Migrated models.")
}
