package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func GetFromCache[T any](ctx context.Context, rkey string, dst *T) error {
	return nil
}

func SaveToCache(ctx context.Context, rkey string, data any) error {
	return nil
}

func DelFromCache(ctx context.Context, rdb *redis.Client, rkeys ...string) error {
	if err := rdb.Del(ctx, rkeys...).Err(); err != nil {
		return err
	}
	return nil
}
