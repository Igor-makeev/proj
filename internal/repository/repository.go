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
}

type Repository struct {
	Authorization
	LoyaltyServiceStorage
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: postgress.NewAuthPostgress(db),
	}

}
