package repository

import (
	"context"
	"log"
	"strconv"

	"github.com/ewallet-backend/internal/dto"
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
	sql := `SELECT id, email, COALESCE(fullname,''), COALESCE(photo_path,''), COALESCE(phone_number,''), created_at, updated_at FROM users WHERE id = $1;`
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
    WHEN t.type = 'TRANSFER_IN' AND t.status = 'SUCCESS'
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

func (u *UserRepository) EditProfile(ctx context.Context, id int, fullname, phone, photo *string) (model.User, error) {
	sql := `UPDATE users 
    SET 
        fullname = $2,
        phone_number = $3,
        photo_path = $4,
        updated_at = NOW()
    WHERE id = $1
    RETURNING fullname, email, phone_number, photo_path`

	var user model.User
	if err := u.db.QueryRow(ctx, sql, id, fullname, phone, photo).Scan(
		&user.Fullname, &user.Email, &user.Phone_number, &user.Photo_path,
	); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *UserRepository) EditUserPin(ctx context.Context, tokenid int, hasHedpin *string) error {
	sql := `UPDATE users SET pin = $2, updated_at = NOW() WHERE id = $1;`
	args := []any{tokenid, hasHedpin}

	_, err := u.db.Exec(ctx, sql, args...)
	return err
}

func (u *UserRepository) EditPassword(ctx context.Context, id int, hashedPassword *string) error {
	sql := `UPDATE users SET password = $2, updated_at = NOW() WHERE id = $1;`
	args := []any{id, hashedPassword}

	_, err := u.db.Exec(ctx, sql, args...)
	return err
}

func (u *UserRepository) CheckPin(ctx context.Context, id int) (*string, error) {
	sql := `SELECT pin FROM users WHERE id = $1;`
	var user *string
	if err := u.db.QueryRow(ctx, sql, id).Scan(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserRepository) GetTransactionReport(ctx context.Context, id int, timePeriod string) ([]model.TransactionReport, error) {
	sql := `SELECT  DATE_TRUNC($2, created_at)::DATE AS period, COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income, COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense
	FROM transactions
	WHERE user_id = $1 AND status = 'SUCCESS'
	GROUP BY DATE_TRUNC($2, created_at)
	ORDER BY period ASC;
`
	args := []any{id, timePeriod}
	var data []model.TransactionReport
	rows, err := u.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.TransactionReport
		if err := rows.Scan(&transaction.Period, &transaction.Income, &transaction.Expense); err != nil {
			return nil, err
		}
		data = append(data, transaction)
	}
	return data, nil
}

func (r *UserRepository) GetTransactionHistory(ctx context.Context, id int, req dto.TransactionHistoryRequest) ([]model.TransactionHistory, error) {
	limit := 10
	page := 1
	if req.Page != "" {
		if p, _ := strconv.Atoi(req.Page); p > 0 {
			page = p
		}
	}
	offset := (page - 1) * limit
	search := req.Search

	query := `
		SELECT 
			t.id AS transaction_id,
			t.amount,
			t.type,
			t.status,
			t.created_at,
			td.description AS transfer_description,
			u_receiver.fullname AS receiver_name,
			pm.name AS payment_method_name,
			COUNT(*) OVER() AS total_count
		FROM transactions t
		LEFT JOIN transfer_details td ON t.id = td.transaction_id
		LEFT JOIN wallet w_receiver ON td.receiver_wallet_id = w_receiver.id
		LEFT JOIN users u_receiver ON w_receiver.user_id = u_receiver.id
		LEFT JOIN topup_details tp ON t.id = tp.transaction_id
		LEFT JOIN payment_methods pm ON tp.payment_method_id = pm.id
		WHERE t.user_id = $1
		  AND (
			$2 = '' OR
			u_receiver.fullname ILIKE '%' || $2 || '%' OR
			pm.name ILIKE '%' || $2 || '%' OR
			td.description ILIKE '%' || $2 || '%'
		  )
		ORDER BY t.created_at DESC
		LIMIT $3 OFFSET $4`

	args := []any{id, search, limit, offset}

	var transactions []model.TransactionHistory
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction model.TransactionHistory
		if err := rows.Scan(
			&transaction.TransactionID,
			&transaction.Amount,
			&transaction.Flow_type,
			&transaction.Type,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.Description,       // td.description
			&transaction.ReceiverName,      // u_receiver.full_name
			&transaction.PaymentMethodName, // pm.name
			&transaction.TotalCount,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	log.Println("apakah repo: ?", transactions)
	return transactions, nil
}
