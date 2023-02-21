package postgress

import (
	"context"
	"proj/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func NewPostgresClient(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		logrus.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	_, err = conn.Exec(context.Background(), usersTableSchema)
	logrus.Print(err)
	_, err = conn.Exec(context.Background(), Index)
	logrus.Print(err)
	return conn, err
}
