package postgress

import (
	"context"
	"proj/internal/entities/models"
	"proj/internal/entities/myerrors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type StoragePostgress struct {
	db *pgxpool.Pool
}

func NewStoragePostgress(db *pgxpool.Pool) *StoragePostgress {
	return &StoragePostgress{db: db}
}

func (sp *StoragePostgress) SaveOrder(ctx context.Context, order models.OrderDTO) error {
	numberInt, _ := strconv.Atoi(order.Number)
	var timeFromDB time.Time
	timeCreated := time.Now()
	var idFromDb int
	err := sp.db.QueryRow(ctx, "insert into orders_table (user_id,number,status,accrual,uploaded_at) values($1,$2,$3,$4,$5)on conflict (number) do update set number =EXCLUDED.number returning user_id,uploaded_at;", order.UserID, numberInt, order.Status, order.Accrual, order.UploadedAt).Scan(&idFromDb, &timeFromDB)
	if err != nil {
		return err
	}

	if timeCreated.Format(time.StampMilli) != timeFromDB.Format(time.StampMilli) {
		if idFromDb != order.UserID {
			return myerrors.ErrOrdUsrConfl
		}
		return myerrors.ErrOrdOverLap
	}

	return nil
}

func (sp *StoragePostgress) OrderUpdate(ctx context.Context, order models.OrderDTO) {
	numberInt, _ := strconv.Atoi(order.Number)
	_, err := sp.db.Exec(ctx, "update orders_table set status=$1, accrual=$2 where number=$3;", order.Status, order.Accrual, numberInt)
	if err != nil {
		logrus.Println(err)
	}
	if order.Accrual > 0 {
		_, err := sp.db.Exec(ctx, "update users_table set ballance=ballance+$1 where id=$2;", order.Accrual, order.UserID)
		if err != nil {
			logrus.Println(err)
		}
	}

}
