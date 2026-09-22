package booking

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBooking(ctx context.Context, userID, roomID string, startTimeStr, endTimeStr string) (string, float64, error) {
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return "", 0, errors.New("invalid start_time format")
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return "", 0, errors.New("invalid end_time format")
	}

	if startTime.Before(time.Now()) {
		return "", 0, errors.New("cannot book in the past")
	}
	if !endTime.After(startTime) {
		return "", 0, errors.New("end_time must be after start_time")
	}

	durationHours := endTime.Sub(startTime).Hours()

	pricePerHour, err := s.repo.GetRoomPricePerHour(ctx, roomID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get room price: %w", err)
	}

	totalPrice := pricePerHour * durationHours

	bookingID, err := s.repo.CreateBookingSafe(ctx, userID, roomID, startTime, endTime, totalPrice)
	if err != nil {
		return "", 0, err
	}

	return bookingID, totalPrice, nil
}

func (s *Service) GetMyBookings(ctx context.Context, userID string) ([]Booking, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) CancelBooking(ctx context.Context, bookingID, userID string) error {
	return s.repo.CancelBooking(ctx, bookingID, userID)
}