package model

type CashFlow struct {
	Balance int `db:"balance"`
	Income  int `db:"income"`
	Expense int `db:"expense"`
}
