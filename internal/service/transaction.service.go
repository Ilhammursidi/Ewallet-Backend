package service

import (
	"context"
	"fmt"
	"math"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TransactionService struct {
	db                    *pgxpool.Pool
	transactionRepository *repository.TransactionRepository
	rdb                   *redis.Client
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, db *pgxpool.Pool, rdb *redis.Client) *TransactionService {
	return &TransactionService{
		db:                    db,
		transactionRepository: transactionRepository,
		rdb:                   rdb,
	}
}

func (ts *TransactionService) FindReceivers(ctx context.Context, userId int, search string, page, limit int) (dto.ReceiverListResponse, error) {
	offset := (page - 1) * limit

	totalId, err := ts.transactionRepository.GetAllReceivers(ctx, ts.db, userId)
	totalPage := int(math.Ceil(float64(totalId) / float64(limit)))

	receivers, err := ts.transactionRepository.FindReceivers(ctx, ts.db, userId, search, limit, offset)
	if err != nil {
		return dto.ReceiverListResponse{}, err
	}

	items := make([]dto.ReceiverResponse, 0, len(receivers))

	for _, receiver := range receivers {

		items = append(items, dto.ReceiverResponse{
			Id:          receiver.Id,
			Photo:       receiver.Photo,
			Receiver:    receiver.Receiver,
			PhoneNumber: receiver.PhoneNumber,
		})
	}
	var NextPage string
	Domain := "http://localhost:8081"
	if page < totalPage {
		Next := page + 1
		NextPage = fmt.Sprintf("%s/transaction/receivers?page=%d&limit=%d", Domain, Next, limit)
	} else {
		NextPage = ""
	}

	var PrevPage string
	if page > totalPage {
		Prev := page - 1
		PrevPage = fmt.Sprintf("%s/transaction/receivers?page=%d&limit=%d", Domain, Prev, limit)
	} else {
		PrevPage = ""
	}

	return dto.ReceiverListResponse{
		Items: items,
		Meta: dto.PaginationMetaResponse{
			Page:       page,
			Limit:      limit,
			Total_Data: totalId,
			Total_Page: totalPage,
			Next_Page:  NextPage,
			Prev_Page:  PrevPage,
		},
	}, nil
}

func (ts *TransactionService) TopUp(ctx context.Context, req dto.TopUpServiceRequest) (*dto.TopUpResponse, error) {
	totalAmount := req.OrderAmount + req.TaxAmount + req.DeliveryFee
	creditAmount := req.OrderAmount

	if creditAmount <= 0 {
		return nil, fmt.Errorf("order_amount too small to cover tax and delivery fee")
	}

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	walletID, err := ts.transactionRepository.GetWalletByUserID(ctx, tx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("TopUp - get wallet: %w", err)
	}

	transactionID, err := ts.transactionRepository.CreateTransaction(ctx, tx, dto.CreateTransactionParams{
		ReceiverWalletID: walletID,
		PaymentMethodID:  req.PaymentMethodID,
		Amount:           totalAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("TopUp - create transaction: %w", err)
	}

	detailID, err := ts.transactionRepository.CreateTopUpDetail(ctx, tx, dto.CreateTopUpDetailParams{
		TransactionID:   transactionID,
		WalletID:        walletID,
		PaymentMethodID: req.PaymentMethodID,
		OrderAmount:     req.OrderAmount,
		TaxAmount:       req.TaxAmount,
		DeliveryFee:     req.DeliveryFee,
		TotalAmount:     totalAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("TopUp - create detail: %w", err)
	}

	if err := ts.transactionRepository.CreditWallet(ctx, tx, walletID, creditAmount); err != nil {
		return nil, fmt.Errorf("TopUp - credit wallet: %w", err)
	}

	if err := ts.transactionRepository.UpdateTopUpStatus(ctx, tx, transactionID, dto.StatusSuccess); err != nil {
		return nil, fmt.Errorf("TopUp - update status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("TopUp - commit: %w", err)
	}

	return &dto.TopUpResponse{
		TransactionID: transactionID,
		TopUpDetailID: detailID,
		TotalAmount:   totalAmount,
		CreditAmount:  creditAmount,
		Status:        dto.StatusSuccess,
	}, nil
}

func (ts *TransactionService) Transfer(ctx context.Context, req dto.TransferServiceRequest) (*dto.TransferResponse, error) {
	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	senderWalletID, err := ts.transactionRepository.GetWalletByUserID(ctx, tx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("Transfer - get sender wallet: %w", err)
	}

	receiverWalletID, err := ts.transactionRepository.GetWalletByUserID(ctx, tx, req.ReceiverID)
	if err != nil {
		return nil, fmt.Errorf("Transfer - get receiver wallet: %w", err)
	}

	if senderWalletID == receiverWalletID {
		return nil, fmt.Errorf("cannot transfer to yourself")
	}

	balance, err := ts.transactionRepository.GetWalletBalance(ctx, tx, senderWalletID)
	if err != nil {
		return nil, fmt.Errorf("Transfer - get balance: %w", err)
	}
	if balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: have %d, need %d", balance, req.Amount)
	}

	transferOutID, err := ts.transactionRepository.CreateTransfer(ctx, tx, dto.CreateTransferParams{
		UserID:           req.UserID,
		SenderWalletID:   senderWalletID,
		ReceiverWalletID: receiverWalletID,
		Amount:           req.Amount,
	})
	if err != nil {
		return nil, fmt.Errorf("Transfer - create transfer: %w", err)
	}

	if err := ts.transactionRepository.CreateTransferDetail(ctx, tx, transferOutID, senderWalletID, receiverWalletID); err != nil {
		return nil, fmt.Errorf("Transfer - create detail: %w", err)
	}

	if err := ts.transactionRepository.DebitWallet(ctx, tx, senderWalletID, req.Amount); err != nil {
		return nil, fmt.Errorf("Transfer - debit wallet: %w", err)
	}
	if err := ts.transactionRepository.CreditWallet(ctx, tx, receiverWalletID, req.Amount); err != nil {
		return nil, fmt.Errorf("Transfer - credit wallet: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Transfer - commit: %w", err)
	}

	return &dto.TransferResponse{
		TransactionID: transferOutID,
		Amount:        req.Amount,
		Status:        dto.StatusSuccess,
	}, nil
}
