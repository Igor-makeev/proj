package repository

import (
	"context"
	"proj/internal/entities/models"
	"proj/internal/repository/postgress"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, login, password string) (int, error)
}

type LoyaltyServiceStorage interface {
	SaveOrder(ctx context.Context, order models.OrderDTO) error
	OrderUpdate(ctx context.Context, order models.OrderDTO)
	GetOrders(ctx context.Context, id int) ([]models.OrderDTO, error)
	GetBalance(ctx context.Context, id int) (*models.UserBallance, error)
	Withdraw(ctx context.Context, withdraw models.Withdrawal, id int) error
	GetWithdrawals(ctx context.Context, id int) ([]models.Withdrawal, error)
}

type Repository struct {
	Authorization
	LoyaltyServiceStorage
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization:         postgress.NewAuthPostgress(db),
		LoyaltyServiceStorage: postgress.NewStoragePostgress(db),
	}

}
