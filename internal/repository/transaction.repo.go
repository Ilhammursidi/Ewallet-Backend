package repository

import (
	"context"
	"fmt"

	"github.com/ewallet-backend/internal/dto"
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
			LOWER(COALESCE(fullname, '')) LIKE $2
			OR LOWER(email) LIKE $2 
			OR LOWER(COALESCE(phone_number, '')) LIKE $2 
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

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, dbtx DBTX, p dto.CreateTransactionParams) (int, error) {
	sql := `
        INSERT INTO transactions
            (receiver_wallet_id, payment_method_id, type, amount, status)
        VALUES
            ($1, $2, 'TOPUP', $3, 'PENDING')
        RETURNING id`

	var transactionID int
	err := dbtx.QueryRow(ctx, sql,
		p.ReceiverWalletID,
		p.PaymentMethodID,
		p.Amount,
	).Scan(&transactionID)
	if err != nil {
		return 0, fmt.Errorf("CreateTransaction: %w", err)
	}
	return transactionID, nil
}

func (tr *TransactionRepository) CreateTopUpDetail(ctx context.Context, dbtx DBTX, p dto.CreateTopUpDetailParams) (int, error) {
	sql := `
        INSERT INTO topup_details
            (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status)
        VALUES
            ($1, $2, $3, $4, $5, $6, $7, 'PENDING')
        RETURNING id`

	var detailID int
	err := dbtx.QueryRow(ctx, sql,
		p.TransactionID,
		p.WalletID,
		p.PaymentMethodID,
		p.OrderAmount,
		p.TaxAmount,
		p.DeliveryFee,
		p.TotalAmount,
	).Scan(&detailID)
	if err != nil {
		return 0, fmt.Errorf("CreateTopUpDetail: %w", err)
	}
	return detailID, nil
}

func (tr *TransactionRepository) UpdateTopUpStatus(ctx context.Context, dbtx DBTX, transactionID int, status dto.TransactionStatus) error {
	sqlTx := `UPDATE transactions SET status = $1 WHERE id = $2`
	sqlDetail := `UPDATE topup_details SET status = $1 WHERE transaction_id = $2`

	if _, err := dbtx.Exec(ctx, sqlTx, status, transactionID); err != nil {
		return fmt.Errorf("UpdateTopUpStatus (transactions): %w", err)
	}
	if _, err := dbtx.Exec(ctx, sqlDetail, status, transactionID); err != nil {
		return fmt.Errorf("UpdateTopUpStatus (topup_details): %w", err)
	}
	return nil
}

func (tr *TransactionRepository) CreditWallet(ctx context.Context, dbtx DBTX, walletID, amount int) error {
	sql := `UPDATE wallet SET balance = balance + $1, updated_at = NOW() WHERE id = $2`
	if _, err := dbtx.Exec(ctx, sql, amount, walletID); err != nil {
		return fmt.Errorf("CreditWallet: %w", err)
	}
	return nil
}

func (tr *TransactionRepository) CreateTransfer(ctx context.Context, dbtx DBTX, p dto.CreateTransferParams) (int, error) {
	sql := `
        INSERT INTO transactions
            (sender_wallet_id, receiver_wallet_id, type, amount, status)
        VALUES
            ($1, $2, 'TRANSFER', $3, 'SUCCESS')
        RETURNING id`

	var transactionID int
	err := dbtx.QueryRow(ctx, sql,
		p.SenderWalletID,
		p.ReceiverWalletID,
		p.Amount,
	).Scan(&transactionID)
	if err != nil {
		return 0, fmt.Errorf("CreateTransferOut: %w", err)
	}
	return transactionID, nil
}

func (tr *TransactionRepository) CreateTransferDetail(ctx context.Context, dbtx DBTX, transactionID, senderWalletID, receiverWalletID int) error {
	sql := `
		INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id)
		VALUES ($1, $2, $3)`

	if _, err := dbtx.Exec(ctx, sql, transactionID, senderWalletID, receiverWalletID); err != nil {
		return fmt.Errorf("CreateTransferDetail: %w", err)
	}
	return nil
}

func (tr *TransactionRepository) DebitWallet(ctx context.Context, dbtx DBTX, walletID, amount int) error {
	sql := `UPDATE wallet SET balance = balance - $1, updated_at = NOW() WHERE id = $2`
	if _, err := dbtx.Exec(ctx, sql, amount, walletID); err != nil {
		return fmt.Errorf("DebitWallet: %w", err)
	}
	return nil
}

func (tr *TransactionRepository) GetWalletBalance(ctx context.Context, dbtx DBTX, walletID int) (int, error) {
	sql := `SELECT balance FROM wallet WHERE id = $1`
	var balance int
	if err := dbtx.QueryRow(ctx, sql, walletID).Scan(&balance); err != nil {
		return 0, fmt.Errorf("GetWalletBalance: %w", err)
	}
	return balance, nil
}

func (tr *TransactionRepository) GetUserIDByWalletID(ctx context.Context, dbtx DBTX, walletID int) (int, error) {
	sql := `SELECT user_id FROM wallet WHERE id = $1`
	var userID int
	if err := dbtx.QueryRow(ctx, sql, walletID).Scan(&userID); err != nil {
		return 0, fmt.Errorf("GetUserIDByWalletID: %w", err)
	}
	return userID, nil
}

func (tr *TransactionRepository) UpdateTransferStatus(ctx context.Context, dbtx DBTX, transactionID int, status dto.TransactionStatus) error {
	sql := `UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2`
	if _, err := dbtx.Exec(ctx, sql, status, transactionID); err != nil {
		return fmt.Errorf("UpdateTransferStatus: %w", err)
	}
	return nil
}

func (tr *TransactionRepository) GetWalletByUserID(ctx context.Context, dbtx DBTX, userID int) (int, error) {
	sql := `SELECT id FROM wallet WHERE user_id = $1`
	var walletID int
	if err := dbtx.QueryRow(ctx, sql, userID).Scan(&walletID); err != nil {
		return 0, fmt.Errorf("GetWalletByUserID: %w", err)
	}
	return walletID, nil
}
