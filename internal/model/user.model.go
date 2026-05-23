package model

type CashFlow struct {
	Balance int `db:"balance"`
	Income  int `db:"income"`
	Expense int `db:"expense"`
}

type UserPin struct {
	Pin string `db:"pin"`
}

type TransactionReport struct {
	Period  string `db:"day"`
	Income  int    `db:"income"`
	Expense int    `db:"expense"`
}
