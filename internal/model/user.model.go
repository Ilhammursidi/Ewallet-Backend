package model

import "time"

type CashFlow struct {
	Balance int `db:"balance"`
	Income  int `db:"income"`
	Expense int `db:"expense"`
}

type UserPin struct {
	Pin string `db:"pin"`
}

type TransactionReport struct {
	Period  *time.Time `db:"day"`
	Income  int        `db:"income"`
	Expense int        `db:"expense"`
}

type TransactionHistory struct {
	TransactionID     int       `db:"transaction_id"`
	Amount            int       `db:"amount"`
	Flow_type         string    `db:"type" binding:"income expense"`
	Type              string    `db:"activity_type" binding:"transfer topup"`
	Status            string    `db:"status" binding:"pending success failed"`
	CreatedAt         time.Time `db:"created_at"`
	Description       string    `db:"transfer_description"`
	ReceiverName      string    `db:"receiver_name"`
	PaymentMethodName string    `db:"payment_method_name"`
	TotalCount        int       `db:"total_count"`
}

type UserLogin struct {
	Id           int        `db:"id"`
	Email        string     `db:"email"`
	Password     string     `db:"password"`
	Fullname     *string    `db:"fullname"`
	Pin          *string    `db:"pin"`
	Photo_path   *string    `db:"photo"`
	Phone_number *string    `db:"phone"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdateAt     *time.Time `db:"updated_at"`
}
