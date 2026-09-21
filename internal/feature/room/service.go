package room

import (
	"context"
	"errors"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateRoom(ctx context.Context, spaceID, name string, capacity int, pricePerHour float64) (string, error) {
	if name == "" || capacity <= 0 || pricePerHour <= 0 {
		return "", errors.New("name, capacity and pricePerHour are required")
	}

	id, err := s.repo.Create(ctx, spaceID, name, capacity, pricePerHour)

	if err != nil {
		return "", fmt.Errorf("failed to create room in db: %w", err)
	}

	return id, nil
}

func (s *Service) GetAllRooms(ctx context.Context, spaceID string) ([]Room, error) {
	if spaceID == "" {
		return nil, errors.New("space_id is required")
	}

	rooms, err := s.repo.GetAllBySpaceID(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %w", err)
	}

	return rooms, nil
}
