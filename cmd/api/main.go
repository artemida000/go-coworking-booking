package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/artemida000/go-coworking-booking/internal/core/config"
	"github.com/artemida000/go-coworking-booking/internal/core/db"
	"github.com/artemida000/go-coworking-booking/internal/core/logger"
	"github.com/artemida000/go-coworking-booking/internal/core/middleware"
	"github.com/artemida000/go-coworking-booking/internal/feature/space"
	"github.com/artemida000/go-coworking-booking/internal/feature/user"
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
		slog.String("port", cfg.Port))

	// подключение к бд
	dbPool, err := db.New(cfg)
	if err != nil {
		log.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	defer dbPool.Close()

	log.Info("Successfully connected to PostgreSQL!")

	userRepo := user.NewRepository(dbPool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, cfg.JWTSecret)

	spaceRepo := space.NewRepository(dbPool)
	spaceService := space.NewService(spaceRepo)
	spaceHandler := space.NewHandler(spaceService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("POST /api/login", userHandler.Login)

	authMiddleware := middleware.Auth(cfg.JWTSecret)

	protectedSpaceCreate := authMiddleware(
		middleware.RequireRole("admin")(http.HandlerFunc(spaceHandler.Create)),
	)

	mux.Handle("POST /api/spaces", protectedSpaceCreate)
	mux.Handle("GET /api/spaces", authMiddleware(http.HandlerFunc(spaceHandler.GetAll)))

	addr := ":" + cfg.Port
	log.Info("Server is listening", slog.String("addr", addr))

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("Server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
