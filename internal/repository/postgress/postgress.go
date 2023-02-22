package postgress

import (
	"context"
	"proj/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func NewPostgresClient(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	conn, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		logrus.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	if _, err := conn.Exec(context.Background(), usersTableSchema); err != nil {
		logrus.Print(err)
	}

	if _, err := conn.Exec(context.Background(), LoginIndex); err != nil {
		logrus.Print(err)
	}
	if _, err := conn.Exec(context.Background(), ordersTableSchema); err != nil {
		logrus.Print(err)
	}
	if _, err := conn.Exec(context.Background(), OrderIndex); err != nil {
		logrus.Print(err)
	}
	return conn, err
}
