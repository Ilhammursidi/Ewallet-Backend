package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BlacklistRepository struct {
	db *pgxpool.Pool
}

func NewBlacklistRepository(db *pgxpool.Pool) *BlacklistRepository {
	return &BlacklistRepository{
		db: db,
	}
}

func (b *BlacklistRepository) AddToBlackList(ctx context.Context, token string) error {
	sql := `INSERT INTO token_blacklist (token) VALUES ($1)`
	_, err := b.db.Exec(ctx, sql, token)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlacklistRepository) IsBlacklist(ctx context.Context, token string) (bool, error) {
	sql := `SELECT token FROM token_blacklist WHERE token = $1`
	var result string
	err := b.db.QueryRow(ctx, sql, token).Scan(&result)
	if err != nil {
		return false, err
	}
	return true, nil
}
