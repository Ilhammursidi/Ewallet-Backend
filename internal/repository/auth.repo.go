package repository

import (
	"context"
	"errors"
	"log"

	"github.com/ewallet-backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (a *AuthRepository) NewUser(ctx context.Context, email, hashedPwd string) (model.User, error) {
	sql := "WITH register AS (INSERT INTO Users (email,password) VALUES ($1,$2) RETURNING id, email, created_at), create_wallet AS (INSERT INTO wallet (user_id) SELECT id FROM register) SELECT id, email, created_at FROM register"

	var user model.User
	if err := a.db.QueryRow(ctx, sql, email, hashedPwd).Scan(&user.Id, &user.Email, &user.Created_at); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (a *AuthRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `SELECT id, password FROM users WHERE email = $1`

	var user model.User
	if err := a.db.QueryRow(ctx, sql, email).Scan(
		&user.Id,
		&user.Password,
	); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (a *AuthRepository) UpdatePin(ctx context.Context, userId int, hashedPin string) error {
	sql := `UPDATE users SET pin = $1 WHERE id = $2`
	_, err := a.db.Exec(ctx, sql, hashedPin, userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) GetUserById(ctx context.Context, userId int) (model.User, error) {
	sql := `SELECT id, pin FROM users WHERE id = $1`
	var user model.User
	if err := a.db.QueryRow(ctx, sql, userId).Scan(&user.Id, &user.Pin); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (a *AuthRepository) AddToBlackList(ctx context.Context, token string) error {
	sql := `INSERT INTO token_blacklist (token) VALUES ($1)`
	log.Println("repository")
	_, err := a.db.Exec(ctx, sql, token)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) IsBlacklist(ctx context.Context, token string) (bool, error) {
	sql := `SELECT token FROM token_blacklist WHERE token = $1;`
	var result string
	err := a.db.QueryRow(ctx, sql, token).Scan(&result)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return true, err
	}
	return true, nil
}
