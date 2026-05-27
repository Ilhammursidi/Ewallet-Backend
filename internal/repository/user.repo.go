package repository

import (
	"context"
	"fmt"
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
	sql := `
        SELECT 
            w.balance AS balance,
            SUM(
                CASE
                    WHEN t.type = 'TRANSFER_IN' AND t.status = 'SUCCESS'
                    THEN t.amount
                    ELSE 0
                END
            ) AS income,
            SUM(
                CASE
                    WHEN t.type = 'TRANSFER_OUT' AND t.status = 'SUCCESS'
                    THEN t.amount
                    ELSE 0
                END
            ) AS expense
        FROM transactions t
        JOIN wallet w ON w.user_id = t.user_id
        WHERE t.user_id = $1
        GROUP BY w.balance`

	var money model.CashFlow
	if err := u.db.QueryRow(ctx, sql, id).Scan(
		&money.Balance,
		&money.Income,
		&money.Expense,
	); err != nil {
		return model.CashFlow{}, err
	}
	return money, nil
}

func (u *UserRepository) EditProfile(ctx context.Context, id int, fullname, phone, photo *string) (model.User, error) {
	sql := `UPDATE users 
    SET 
        fullname = COALESCE($2, fullname),
        phone_number = COALESCE($3, phone_number),
        photo_path = COALESCE($4, photo_path),
        updated_at = NOW()
    WHERE id = $1
	RETURNING id, fullname, phone_number, photo_path;`

	args := []any{id, fullname, phone, photo}

	var user model.User
	err := u.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Fullname, &user.Phone_number, &user.Photo_path)
	if err != nil {
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
func (u *UserRepository) GetTransactionReport(ctx context.Context, id int, timePeriod string) ([]dto.TransactionReportDTO, error) {
	var sql string

	switch timePeriod {
	case "week":
		sql = `
			WITH date_series AS (
				SELECT generate_series(
					(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')::DATE - INTERVAL '6 days',
					(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')::DATE,
					INTERVAL '1 day'
				)::DATE AS period
			)
			SELECT 
				ds.period,
				COALESCE(SUM(CASE WHEN t.type = 'TRANSFER_IN' AND t.status = 'SUCCESS' THEN t.amount ELSE 0 END), 0) AS total_income,
				COALESCE(SUM(CASE WHEN t.type = 'TRANSFER_OUT' AND t.status = 'SUCCESS' THEN t.amount ELSE 0 END), 0) AS total_expense
			FROM date_series ds
			LEFT JOIN transactions t
				ON (t.created_at AT TIME ZONE 'Asia/Jakarta')::DATE = ds.period
				AND t.user_id = $1
			GROUP BY ds.period
			ORDER BY ds.period ASC`

	case "month":
		sql = `
			WITH date_series AS (
				SELECT generate_series(
					DATE_TRUNC('month', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')::DATE,
					(DATE_TRUNC('month', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta') + INTERVAL '1 month - 1 day')::DATE,
					INTERVAL '1 day'
				)::DATE AS period
			)
			SELECT 
				ds.period,
				COALESCE(SUM(CASE WHEN t.type = 'TRANSFER_IN' AND t.status = 'SUCCESS' THEN t.amount ELSE 0 END), 0) AS total_income,
				COALESCE(SUM(CASE WHEN t.type = 'TRANSFER_OUT' AND t.status = 'SUCCESS' THEN t.amount ELSE 0 END), 0) AS total_expense
			FROM date_series ds
			LEFT JOIN transactions t
				ON (t.created_at AT TIME ZONE 'Asia/Jakarta')::DATE = ds.period
				AND t.user_id = $1
			GROUP BY ds.period
			ORDER BY ds.period ASC`

	case "year":
		sql = `
			WITH date_series AS (
				SELECT generate_series(
					DATE_TRUNC('year', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')::DATE,
					(DATE_TRUNC('year', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta') + INTERVAL '11 months')::DATE,
					INTERVAL '1 month'
				)::DATE AS period
			)
			SELECT 
				ds.period,
				COALESCE(SUM(CASE WHEN t.type = 'TRANSFER_IN' AND t.status = 'SUCCESS' THEN t.amount ELSE 0 END), 0) AS total_income,
				COALESCE(SUM(CASE WHEN t.type = 'TRANSFER_OUT' AND t.status = 'SUCCESS' THEN t.amount ELSE 0 END), 0) AS total_expense
			FROM date_series ds
			LEFT JOIN transactions t
				ON DATE_TRUNC('month', t.created_at AT TIME ZONE 'Asia/Jakarta')::DATE = ds.period
				AND t.user_id = $1
			GROUP BY ds.period
			ORDER BY ds.period ASC`

	default:
		return nil, fmt.Errorf("invalid period: %s (week/month/year)", timePeriod)
	}

	rows, err := u.db.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("GetTransactionReport: %w", err)
	}
	defer rows.Close()

	var data []dto.TransactionReportDTO
	for rows.Next() {
		var report dto.TransactionReportDTO
		if err := rows.Scan(&report.Period, &report.Income, &report.Expense); err != nil {
			return nil, fmt.Errorf("GetTransactionReport scan: %w", err)
		}
		data = append(data, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTransactionReport rows: %w", err)
	}

	return data, nil
}

func (r *UserRepository) GetTransactionHistory(ctx context.Context, id int, req dto.TransactionHistoryRequest) ([]dto.TransactionHistoryDTO, error) {
	limit := 10
	page := 1
	if req.Page != "" {
		if p, _ := strconv.Atoi(req.Page); p > 0 {
			page = p
		}
	}
	offset := (page - 1) * limit

	query := `
		SELECT 
			t.id                                AS transaction_id,
			t.amount,
			t.type,
			t.flow_type,
			t.status,
			t.created_at,
			COALESCE(pm.payment_name, '')       AS payment_method_name,
			w_receiver.user_id                  AS receiver_id,
			COALESCE(u_receiver.fullname, '')   AS receiver_name,
			COALESCE(u_sender.fullname, '')     AS sender_name,
			COUNT(*) OVER()                     AS total_count
		FROM transactions t
		LEFT JOIN topup_details tp      ON t.id = tp.transaction_id
		LEFT JOIN payment_methods pm    ON tp.payment_method_id = pm.id
		LEFT JOIN transfer_details trd  ON t.id = trd.transaction_id
		LEFT JOIN wallet w_receiver     ON trd.receiver_wallet_id = w_receiver.id
		LEFT JOIN users u_receiver      ON w_receiver.user_id = u_receiver.id
		LEFT JOIN wallet w_sender       ON trd.sender_wallet_id = w_sender.id
		LEFT JOIN users u_sender        ON w_sender.user_id = u_sender.id
		WHERE t.user_id = $1
		  AND (
			$2 = '' OR
			u_receiver.fullname ILIKE '%' || $2 || '%' OR
			u_sender.fullname   ILIKE '%' || $2 || '%' OR
			pm.payment_name     ILIKE '%' || $2 || '%' OR
			t.type::TEXT        ILIKE '%' || $2 || '%'
		  )
		ORDER BY t.created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(ctx, query, id, req.Search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetTransactionHistory query: %w", err)
	}
	defer rows.Close()

	var data []dto.TransactionHistoryDTO
	for rows.Next() {
		var t dto.TransactionHistoryDTO
		if err := rows.Scan(
			&t.TransactionID,
			&t.Amount,
			&t.Type,
			&t.FlowType,
			&t.Status,
			&t.CreatedAt,
			&t.PaymentMethodName,
			&t.ReceiverID,
			&t.ReceiverName,
			&t.SenderName,
			&t.TotalCount,
		); err != nil {
			log.Println("GetTransactionHistory scan error:", err)
			return nil, fmt.Errorf("GetTransactionHistory scan: %w", err)
		}
		data = append(data, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTransactionHistory rows: %w", err)
	}

	return data, nil
}
