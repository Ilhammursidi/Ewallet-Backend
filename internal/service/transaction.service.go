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
