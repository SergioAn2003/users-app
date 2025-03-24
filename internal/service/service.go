package service

import (
	"context"
	"users-app/internal/entity"

	"github.com/gofrs/uuid/v5"
)

//go:generate go run go.uber.org/mock/mockgen@latest -source=service.go -destination=../mocks/service.go -package=mocks -typed

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *Service) CreateUser(ctx context.Context, user entity.User) error {
	return s.userRepo.CreateUser(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, user entity.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}
