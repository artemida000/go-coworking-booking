package space

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

func (s *Service) CreateSpace(ctx context.Context, name, address, description string) (string, error) {
	if name == "" || address == "" {
		return "", errors.New("name and address are required")
	}

	id, err := s.repo.Create(ctx, name, address, description)
	
	if err != nil {
		return "", fmt.Errorf("failed to create space in db: %w", err)
	}

	return id, nil
}

func (s *Service) GetAllSpace(ctx context.Context) ([]Space, error) {
	spaces, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get spaces fron db: %w", err)
	}

	return spaces, nil
}