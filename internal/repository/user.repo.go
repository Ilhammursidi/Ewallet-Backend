package repository

import (
	"context"
	"log"

	"github.com/ewallet-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetProfileId(ctx context.Context, id int) (model.User, error) {
	sql := `SELECT id, email, fullname, photo_path, phone_number, created_at, updated_at FROM users WHERE id = $1;`
	args := []any{id}

	var user model.User
	if err := u.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.Fullname, &user.Photo_path, &user.Phone_number, &user.Created_at, &user.Updated_at); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *UserRepository) GetMoneyAccountInfo(ctx context.Context, id int) (model.CashFlow, error) {
	sql := `SELECT w.balance AS balance, SUM(
	CASE
    WHEN t.type IN ('TOPUP','TRANSFER_IN') AND t.status = 'SUCCESS'
	THEN t.amount
	ELSE 0
	END
) AS income, SUM(
	CASE
    WHEN t.type = 'TRANSFER_OUT' AND t.status = 'SUCCESS'
	THEN t.amount
	ELSE 0
	END
) AS expense 
FROM transactions t
JOIN wallet w ON w.id = t.user_id
WHERE t.user_id = $1
GROUP BY w.balance;`

	args := []any{id}

	var money model.CashFlow
	if err := u.db.QueryRow(ctx, sql, args...).Scan(&money.Balance, &money.Expense, &money.Income); err != nil {
		return model.CashFlow{}, err
	}
	return money, nil
}

// func (u *UserRepository) CreateTransfer(ctx context.Context, id int)

func (u *UserRepository) EditProfile(ctx context.Context, id int, fullname, phone, photo *string) (model.User, error) {
	log.Println(*fullname)
	sql := `UPDATE users 
	SET 
		fullname = $2,
		phone_number = $3,
		photo_path = $4,
		updated_at = NOW()
		WHERE id = $1
	RETURNING fullname, email, phone_number, photo_path;
	`
	args := []any{id, fullname, phone, photo}

	var user model.User
	if err := u.db.QueryRow(ctx, sql, args...).Scan(&user.Fullname, &user.Email, &user.Phone_number, &user.Photo_path); err != nil {
		return model.User{}, err
	}
	return user, nil
}
