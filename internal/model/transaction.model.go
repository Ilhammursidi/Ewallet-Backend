package model

type Receiver struct {
	Id          int     `db:"id"`
	Photo       *string `db:"photo"`
	Receiver    *string `db:"receiver"`
	PhoneNumber *string `db:"phone"`
}
