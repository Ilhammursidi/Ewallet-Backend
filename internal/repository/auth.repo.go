package repository

import (
	"context"

	"github.com/ewallet-backend/internal/model"
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

func (a *AuthRepository) FetchUserList(ctx context.Context) ([]model.User, error) {
	sql := "SELECT id,email,password,pin,fullname,photo_path,phone_number,created_at,updated_at,deleted_at FROM users LIMIT 5"
	rows, err := a.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.User
	for rows.Next() {
		var users model.User
		if err := rows.Scan(&users.Id, &users.Email, &users.Fullname, &users.Password, &users.Pin, &users.Photo_path, &users.Phone_number, &users.Created_at, &users.Updated_at, &users.Deleted_at); err != nil {
			return nil, err
		}
		result = append(result, users)
	}
	return result, nil
}

func (a *AuthRepository) NewUser(ctx context.Context, email, hashedPwd string) (model.User, error) {
	sql := "WITH register AS (INSERT INTO Users (email,password) VALUES ($1,$2) RETURNING id, email, created_at), create_wallet AS (INSERT INTO wallet (user_id) SELECT id FROM register) SELECT id, email, created_at FROM register"
	args := []any{email, hashedPwd}

	var user model.User
	if err := a.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.Created_at); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (a *AuthRepository) GetUserById(ctx context.Context, id int) (model.User, error) {
	sql := `SELECT id, email, created_at FROM users WHERE id = $1`

	var user model.User
	if err := a.db.QueryRow(ctx, sql, id).Scan(
		&user.Id,
		&user.Email,
		&user.Created_at,
	); err != nil {
		return model.User{}, err
	}
	return user, nil
}
