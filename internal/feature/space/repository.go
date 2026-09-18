package space

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Space struct {
	ID          string
	Name        string
	Address     string
	Description string
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, name, address, description string) (string, error) {
	query := `
		INSERT INTO spaces (name, address, description)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id string

	err := r.db.QueryRow(ctx, query, name, address, description).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}