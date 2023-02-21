package postgress

import (
	"context"
	"proj/internal/entities/models"
	"proj/internal/entities/myerrors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgress(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user models.User) error {

	_, err := r.db.Exec(ctx, "insert into users_table (login,password_hash) values ($1,$2);", user.Login, user.Password)
	if err != nil {
		pqErr := err.(*pgconn.PgError)
		if pqErr.Code == pgerrcode.UniqueViolation {
			return &myerrors.LoginConflict{Elem: user.Login}
		}
		return err
	}

	return nil
}

func (r *AuthPostgres) GetUser(ctx context.Context, login, password string) (int, error) {
	var userID int
	if err := r.db.QueryRow(ctx, "select id from users_table where login=$1 and password_hash=$2", login, password).Scan(&userID); err != nil {

		if err == pgx.ErrNoRows {
			return 0, myerrors.InvalidLoginOrPassword
		}
		return 0, err

	}
	return userID, nil
}
