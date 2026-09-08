package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, email, passwordHash string) (string, error) {
	query := `
		INSERT INTO USERS (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id string 

	err := r.db.QueryRow(ctx, query, email, passwordHash).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}


func (r *Repository) GetEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role
		FROM users
		Where email = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &user, nil
}