package model

import "time"

type User struct {
	Id           int        `db:"id"`
	Email        string     `db:"email" example:"anjay"`
	Password     string     `db:"password"`
	Pin          *string    `db:"pin"`
	Fullname     *string    `db:"fullname"`
	Photo_path   *string    `db:"photo"`
	Phone_number *string    `db:"phone"`
	IsVerified   *bool      `db:"isverified"`
	Created_at   *time.Time `db:"created_at"`
	Updated_at   *time.Time `db:"updated_at"`
	Deleted_at   *time.Time `db:"deleted_at"`
}
