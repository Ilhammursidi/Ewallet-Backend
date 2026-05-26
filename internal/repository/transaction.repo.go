package repository

import (
	"context"

	"github.com/ewallet-backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type TransactionRepository struct{}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{}
}

func (tr *TransactionRepository) GetAllReceivers(ctx context.Context, dbtx DBTX, userId int) (int, error) {
	sql := `SELECT COUNT(id) FROM users WHERE id != $1`
	var count int
	err := dbtx.QueryRow(ctx, sql, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil

}

func (tr *TransactionRepository) FindReceivers(ctx context.Context, dbtx DBTX, userId int, search string, limit, offset int) ([]model.Receiver, error) {
	target := `%` + search + `%`
	sqlQuery := `
		SELECT id, COALESCE(photo_path, '') AS photo, COALESCE(fullname, email) AS receiver, COALESCE(phone_number, '') AS phone
		FROM users
		WHERE id != $1
		AND (
			COALESCE(fullname, '') ILIKE $2
			OR email ILIKE $2 
			OR COALESCE(phone_number, '') ILIKE $2 
			)
		ORDER BY COALESCE(fullname, email) ASC
		LIMIT $3
		OFFSET $4;
	`

	args := []any{userId, target, limit, offset}

	rows, err := dbtx.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receivers := []model.Receiver{}
	for rows.Next() {
		var receiver model.Receiver
		if err := rows.Scan(&receiver.Id, &receiver.Photo, &receiver.Receiver, &receiver.PhoneNumber); err != nil {
			return nil, err
		}

		receivers = append(receivers, receiver)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return receivers, nil
}

// func (tr *TransactionRepository) CreateTopUp(ctx context.Context, dbtx DBTX, id int, amount int)   {
// 	sql := `WITH topup AS (
// 	INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status, created_at) VALUES ($1)`
// }
