package booking

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Booking struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RoomID     string    `json:"room_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRoomPricePerHour(ctx context.Context, roomID string) (float64, error) {
	var price float64
	err := r.db.QueryRow(ctx, "SELECT price_per_hour FROM rooms WHERE id = $1", roomID).Scan(&price)
	return price, err
}

func (r *Repository) CreateBookingSafe(ctx context.Context, userID, roomID string, startTime, endTime time.Time, totalPrice float64) (string, error) {
	query := `
		INSERT INTO bookings (user_id, room_id, start_time, end_time, total_price)
		SELECT $1, $2, $3, $4, $5
		WHERE NOT EXISTS (
			SELECT 1 FROM bookings 
			WHERE room_id = $2 
			AND status = 'active'
			AND start_time < $4 
			AND end_time > $3
		)
		RETURNING id
	`

	var id string
	err := r.db.QueryRow(ctx, query, userID, roomID, startTime, endTime, totalPrice).Scan(&id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", errors.New("room is already booked for this time")
		}
		return "", err
	}
	return id, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) ([]Booking, error) {
	query := `
		SELECT id, user_id, room_id, start_time, end_time, status, total_price 
		FROM bookings 
		WHERE user_id = $1 
		ORDER BY start_time DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.RoomID, &b.StartTime, &b.EndTime, &b.Status, &b.TotalPrice)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if bookings == nil {
		bookings = []Booking{}
	}
	return bookings, nil
}

func (r *Repository) CancelBooking(ctx context.Context, bookingID, userID string) error {
	query := `
		UPDATE bookings 
		SET status = 'cancelled' 
		WHERE id = $1 AND user_id = $2 AND status = 'active'
	`
	
	tag, err := r.db.Exec(ctx, query, bookingID, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("booking not found or cannot be cancelled")
	}

	return nil
}