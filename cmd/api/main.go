package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/artemida000/go-coworking-booking/internal/core/config"
	"github.com/artemida000/go-coworking-booking/internal/core/db"
	"github.com/artemida000/go-coworking-booking/internal/core/logger"
	"github.com/artemida000/go-coworking-booking/internal/core/middleware"
	"github.com/artemida000/go-coworking-booking/internal/feature/booking"
	"github.com/artemida000/go-coworking-booking/internal/feature/room"
	"github.com/artemida000/go-coworking-booking/internal/feature/space"
	"github.com/artemida000/go-coworking-booking/internal/feature/user"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	log := logger.Setup(cfg.Env)

	log.Info("Starting Coworking Booking API...",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.Port))

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

	roomRepo := room.NewRepository(dbPool)
	roomService := room.NewService(roomRepo)
	roomHandler := room.NewHandler(roomService)

	bookingRepo := booking.NewRepository(dbPool)
	bookingService := booking.NewService(bookingRepo)
	bookingHandler := booking.NewHandler(bookingService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("POST /api/login", userHandler.Login)

	authMiddleware := middleware.Auth(cfg.JWTSecret)

	protectedSpaceCreate := authMiddleware(
		middleware.RequireRole("admin")(http.HandlerFunc(spaceHandler.Create)),
	)

	protectedRoomCreate := authMiddleware(
		middleware.RequireRole("admin")(http.HandlerFunc(roomHandler.Create)),
	)

	mux.Handle("POST /api/spaces", protectedSpaceCreate)
	mux.Handle("POST /api/rooms", protectedRoomCreate)
	mux.Handle("GET /api/spaces", authMiddleware(http.HandlerFunc(spaceHandler.GetAll)))
	mux.Handle("GET /api/spaces/{id}/rooms", authMiddleware(http.HandlerFunc(roomHandler.GetAll)))
	mux.Handle("POST /api/bookings", authMiddleware(http.HandlerFunc(bookingHandler.Create)))
	mux.Handle("GET /api/bookings/my", authMiddleware(http.HandlerFunc(bookingHandler.GetMyBookings)))
	mux.Handle("DELETE /api/bookings/{id}", authMiddleware(http.HandlerFunc(bookingHandler.Cancel)))

	addr := ":" + cfg.Port
	log.Info("Server is listening", slog.String("addr", addr))

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("Server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
