package room

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Room struct {
	ID           string
	SpaceID      string
	Name         string
	Capacity     int
	PricePerHour float64
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, spaceID, name string, capacity int, pricePerHour float64) (string, error) {
	query := `
		INSERT INTO rooms (space_id, name, capacity, price_per_hour)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string

	err := r.db.QueryRow(ctx, query, spaceID, name, capacity, pricePerHour).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) GetAllBySpaceID(ctx context.Context, spaceID string) ([]Room, error) {
	query := `
		SELECT id, space_id, name, capacity, price_per_hour 
		FROM rooms 
		WHERE space_id = $1
	`

	rows, err := r.db.Query(ctx, query, spaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var rm Room
		err := rows.Scan(&rm.ID, &rm.SpaceID, &rm.Name, &rm.Capacity, &rm.PricePerHour)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if rooms == nil {
		rooms = []Room{}
	}

	return rooms, nil
}