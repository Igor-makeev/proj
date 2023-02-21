package service

import (
	"context"
	"proj/config"
	"proj/internal/entities/models"
	"proj/internal/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) error
	GenerateToken(ctx context.Context, login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		Authorization: NewAuthService(repo, cfg),
	}
}
