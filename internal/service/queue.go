package service

import (
	"context"
	"proj/internal/entities/models"
)

type Client interface {
}
type TxStatus int

const ()

type Queue struct {
	client Client
	buf    chan models.OrderDTO
}

func NewQueue(ctx context.Context, client Client) *Queue {
	q := &Queue{
		buf:    make(chan models.OrderDTO, 100),
		client: client,
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
			q.distributeRequest(ctx, order)

		case <-ctx.Done():

			return
		}
	}
}

func (tm *Queue) distributeRequest(ctx context.Context, order models.OrderDTO) {
	switch order.Status {
	case models.StatusNew:

	case models.StatusNew:
		//..
	case models.StatusNew:
		//..
	case models.StatusNew:
		//..
	}
}
