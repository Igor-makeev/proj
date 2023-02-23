package service

import (
	"context"
	"proj/internal/entities/models"
	"proj/internal/entities/myerrors"
	"proj/internal/repository"
	"proj/pkg/luhn"
	"time"
)

type LoyaltyService struct {
	repo  *repository.Repository
	queue *Queue
}

func NewLoyaltyService(ctx context.Context, repo *repository.Repository, client *Client) *LoyaltyService {
	return &LoyaltyService{
		repo:  repo,
		queue: NewQueue(ctx, client, repo),
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

func (ls *LoyaltyService) OrderUpdate(ctx context.Context, order models.OrderDTO) {
	ls.repo.OrderUpdate(ctx, order)
}

func (ls *LoyaltyService) GetOrders(ctx context.Context, id int) ([]models.OrderDTO, error) {

	return ls.repo.GetOrders(ctx, id)
}

func (ls *LoyaltyService) GetBalance(ctx context.Context, id int) (*models.UserBallance, error) {
	return ls.repo.GetBalance(ctx, id)
}

func (ls *LoyaltyService) Withdraw(ctx context.Context, withdraw models.Withdrawal, id int) error {
	if !luhn.LuhnValidation(string(withdraw.OrderNumber)) {
		return myerrors.ErrInvalidOrderNumber

	}

	accountState, err := ls.repo.GetBalance(ctx, id)
	if err != nil {
		return err
	}

	if accountState.Current < withdraw.Sum {
		return myerrors.ErrNoMoney
	}

	if err := ls.repo.Withdraw(ctx, withdraw, id); err != nil {
		return err
	}
	return nil
}

func (ls *LoyaltyService) GetWithdrawals(ctx context.Context, id int) ([]models.Withdrawal, error) {

	return ls.repo.GetWithdrawals(ctx, id)
}
