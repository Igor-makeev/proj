package service

import (
	"context"
	"proj/internal/entities/models"
)

type communicator interface {
	DoRequest(orderNumber string, UserID int, out chan models.OrderDTO)
}

type updater interface {
	OrderUpdate(ctx context.Context, order models.OrderDTO)
}

type Queue struct {
	updater
	communicator
	buf chan models.OrderDTO
}

func NewQueue(ctx context.Context, client communicator, updater updater) *Queue {
	q := &Queue{
		buf:          make(chan models.OrderDTO, 100),
		communicator: client,
		updater:      updater,
	}
	q.Run(ctx)
	return q
}

func (q *Queue) Run(ctx context.Context) {
	go q.listen(ctx)
}

func (q *Queue) Close() {
	close(q.buf)
}

func (q *Queue) listen(ctx context.Context) {

	for {

		select {
		case order := <-q.buf:
			q.distribute(ctx, order)

		case <-ctx.Done():

			return
		}
	}
}

func (q *Queue) distribute(ctx context.Context, order models.OrderDTO) {
	switch order.Status {
	case models.StatusNew:
		q.communicator.DoRequest(order.Number, order.UserID, q.buf)
	case models.StatusInvalid:
		q.updater.OrderUpdate(ctx, order)
	case models.StatusProcessing:
		q.communicator.DoRequest(order.Number, order.UserID, q.buf)
	case models.StatusProcessed:
		q.updater.OrderUpdate(ctx, order)
	}
}
