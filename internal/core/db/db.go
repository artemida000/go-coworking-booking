package db

import (
	"context"
	"fmt"
	"time"

	"github.com/artemida000/go-coworking-booking/internal/core/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *config.Confing) (*pgxpool.Pool, error) {

	// формируем строку подключения
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmodel=%s",
	cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslmode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// создаем пул соединения
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// пингуем нашу бд
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil 

}