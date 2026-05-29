package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	SaveResetToken(ctx context.Context, token string, userID int, ttl time.Duration) error
	GetUserIDByToken(ctx context.Context, token string) (int, error)
	DeleteResetToken(ctx context.Context, token string) error
}

type cacheRepository struct {
	rdb *redis.Client
}

func NewCacheRepository(rdb *redis.Client) CacheRepository {
	return &cacheRepository{rdb: rdb}
}

func (r *cacheRepository) SaveResetToken(ctx context.Context, token string, userID int, ttl time.Duration) error {
	return r.rdb.Set(ctx, "reset_token:"+token, userID, ttl).Err()
}

func (r *cacheRepository) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	val, err := r.rdb.Get(ctx, "reset_token:"+token).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *cacheRepository) DeleteResetToken(ctx context.Context, token string) error {
	return r.rdb.Del(ctx, "reset_token:"+token).Err()
}
