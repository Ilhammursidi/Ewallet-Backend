package service

import (
	"context"

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

	return dto.ReceiverListResponse{
		Items: items,
		Meta: dto.PaginationMetaResponse{
			Page:  page,
			Limit: limit,
		},
	}, nil
}
