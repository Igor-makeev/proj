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

type LoyaltyServicer interface {
	SaveOrder(ctx context.Context, number string, id int) error
	GetOrders(ctx context.Context, id int) ([]models.OrderDTO, error)
	GetBalance(ctx context.Context, id int) (*models.UserBallance, error)
	Withdraw(ctx context.Context, withdraw models.Withdrawal, id int) error
	GetWithdrawals(ctx context.Context, id int) ([]models.Withdrawal, error)
}

type Service struct {
	Authorization
	LoyaltyServicer
}

func NewService(ctx context.Context, repo *repository.Repository, cfg *config.Config, client *Client) *Service {

	return &Service{
		Authorization:   NewAuthService(repo, cfg),
		LoyaltyServicer: NewLoyaltyService(ctx, repo, client),
	}
}
