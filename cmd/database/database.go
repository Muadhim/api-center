package database

import (
	"api-center/configs"
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *configs.Config) (db *gorm.DB) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.SSLMode,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	sqlDB, _ := db.DB()
	maxIdleConns := 30
	if cfg.SetMaxIdleConns != "" {
		maxIdleConns, _ = strconv.Atoi(cfg.SetMaxIdleConns)
	}
	MaxOpenConns := 100
	if cfg.SetMaxOpenConns != "" {
		MaxOpenConns, _ = strconv.Atoi(cfg.SetMaxOpenConns)
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(MaxOpenConns)

	slog.Info("Database API CENTER connected")

	return db
}
