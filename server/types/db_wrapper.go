package types

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type DBWrapper struct {
	DB *pgx.Conn
}

func (w *DBWrapper) WithTransaction(handler func(c *fiber.Ctx, tx *pgx.Tx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		// Start a new transaction
		tx, err := w.DB.Begin(ctx)
		if err != nil {
			fmt.Printf("error opening transation\n")
			fmt.Printf("err: %+v\n", err)
			return err
		}

		fmt.Printf("transaction opened\n")

		defer func() {
			if p := recover(); p != nil {
				tx.Rollback(ctx)
				panic(p) // re-throw panic after Rollback
			} else if err != nil {
				tx.Rollback(ctx) // err is non-nil; don't change it
			} else {
				err = tx.Commit(ctx) // err is nil; if Commit returns error update err
			}
		}()

		// Run the handler and catch any errors
		err = handler(c, &tx)

		return err
	}
}
