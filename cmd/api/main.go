package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/artemida000/go-coworking-booking/internal/core/config"
	"github.com/artemida000/go-coworking-booking/internal/core/db"
	"github.com/artemida000/go-coworking-booking/internal/core/logger"
)

func main() {

	// передаем переменные окружения
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// инициализируем логер
	log := logger.Setup(cfg.Env)

	log.Info("Starting Coworking Booking API...", 
	slog.String("env", cfg.Env), 
	slog.String("port", cfg.Port),)

	// подключение к бд
	dbPool, err := db.New(cfg)
	if err != nil {
		log.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	defer dbPool.Close()

	log.Info("Successfully connected to PostgreSQL!")

	fmt.Printf("Configuration loaded successfully! Server will run on port: %s\n", cfg.Port)
}