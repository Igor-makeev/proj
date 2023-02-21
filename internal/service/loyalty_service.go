package service

import (
	"context"
	"proj/internal/entities/models"
	"proj/internal/repository"
	"time"
)

type LoyaltyService struct {
	repo  *repository.Repository
	queue *Queue
}

func NewLoyaltyService(ctx context.Context, repo *repository.Repository, client *Client) *LoyaltyService {
	return &LoyaltyService{
		repo:  repo,
		queue: NewQueue(ctx, client),
	}
}

func (ls *LoyaltyService) SaveOrder(ctx context.Context, number string, id int) error {
	neworder := models.OrderDTO{
		UserID:     id,
		Number:     number,
		Status:     models.StatusNew,
		Accrual:    0,
		UploadedAt: time.Now(),
	}
	err := ls.repo.SaveOrder(ctx, neworder)
	if err != nil {
		return err
	}
	ls.queue.buf <- neworder
	return nil
}
