package postgress

import (
	"context"
	"proj/internal/entities/models"
	"proj/internal/entities/myerrors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
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
	var idFromDB int64
	err := sp.db.QueryRow(ctx, "insert into orders_table (user_id,number,status,accrual,uploaded_at) values($1,$2,$3,$4,$5)on conflict (number) do update set number =EXCLUDED.number returning user_id,uploaded_at;", order.UserID, numberInt, order.Status, order.Accrual, order.UploadedAt).Scan(&idFromDB, &timeFromDB)
	if err != nil {

		return err
	}

	if timeCreated.Format(time.StampMilli) != timeFromDB.Format(time.StampMilli) {
		if idFromDB != int64(order.UserID) {
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

		_, err := sp.db.Exec(ctx, "update users_table set current_ballance=current_ballance+$1 where id=(select user_id from orders_table where number=$2);", order.Accrual, numberInt)
		if err != nil {

			logrus.Println(err)
		}
	}

}
func (sp *StoragePostgress) GetOrders(ctx context.Context, id int) ([]models.OrderDTO, error) {

	rows, err := sp.db.Query(ctx, "select number, status, accrual, uploaded_at from orders_table where user_id=$1;", id)
	if err != nil {

		return nil, err
	}
	var list = make([]models.OrderDTO, 0, 100)
	for rows.Next() {

		var order models.OrderDTO
		var number int

		err := rows.Scan(&number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {

			return nil, err
		}
		order.Number = strconv.Itoa(number)
		list = append(list, order)
	}
	return list, nil
}

func (sp *StoragePostgress) GetBalance(ctx context.Context, id int) (*models.UserBallance, error) {
	var balance models.UserBallance

	err := sp.db.QueryRow(ctx, "select current_ballance, withdrawn from users_table	where id=$1;", id).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (sp *StoragePostgress) Withdraw(ctx context.Context, withdraw models.Withdrawal, id int) error {

	numberInt, _ := strconv.Atoi(withdraw.OrderNumber)
	tx, err := sp.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()
	_, err = sp.db.Exec(ctx, "insert into withdrawal_table (user_id,number,sum,uploaded_at) values($1,$2,$3,$4);", id, numberInt, withdraw.Sum, time.Now())
	if err != nil {
		return err
	}

	_, err = sp.db.Exec(ctx, "update users_table set current_ballance=current_ballance-$1, withdrawn=withdrawn+$1 where id=$2;", withdraw.Sum, id)
	if err != nil {

		return err
	}
	return nil
}

func (sp *StoragePostgress) GetWithdrawals(ctx context.Context, id int) ([]models.Withdrawal, error) {
	rows, err := sp.db.Query(ctx, "select * from withdrawal_table where user_id=$1 order by processed_at;", id)
	if err != nil {

		return nil, err
	}

	var list = make([]models.Withdrawal, 0, 100)
	for rows.Next() {
		var withdraw models.Withdrawal

		err := rows.Scan(&withdraw.OrderNumber, &withdraw.Sum, &withdraw.ProcessedAt)

		if err != nil {

			return nil, err
		}
		list = append(list, withdraw)
	}
	return list, nil
}
