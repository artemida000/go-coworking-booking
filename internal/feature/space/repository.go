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

func (r *Repository) GetAll(ctx context.Context) ([]Space, error) {
	query := `
		SELECT id, name, address, description
		FROM spaces
	`
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	var spaces []Space

	for rows.Next() {
		var s Space

		err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.Description)
		if err != nil {
			return nil, err
		}

		spaces = append(spaces, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if spaces == nil {
		spaces = []Space{}
	}

	return spaces, nil
}