// db/seeder/seeder.go
package seeder

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(db *pgxpool.Pool) error {
	if err := seedPaymentMethods(db); err != nil {
		return err
	}
	return nil
}

func seedPaymentMethods(db *pgxpool.Pool) error {
	sql := `
        INSERT INTO payment_methods (payment_name)
        VALUES
		('Bank Rakyat Indonesia'),
		('Dana'),
		('Bank Central Asia'),
		('Gopay'),
		('Ovo')
        ON CONFLICT DO NOTHING`

	_, err := db.Exec(context.Background(), sql)
	return err
}
