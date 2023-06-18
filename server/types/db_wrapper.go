package types

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBWrapper struct {
	DB *pgxpool.Pool
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
				rbErr := tx.Rollback(ctx) // err is non-nil; don't change it
				fmt.Printf("err: %+v\n", rbErr)
				fmt.Printf("err: %+v\n", err)
				panic(p) // re-throw panic after Rollback
			} else if err != nil {
				rbErr := tx.Rollback(ctx) // err is non-nil; don't change it
				fmt.Printf("err: %+v\n", rbErr)
			} else {
				err = tx.Commit(ctx) // err is nil; if Commit returns error update err
				fmt.Printf("err: %+v\n", err)
			}
		}()

		// Run the handler and catch any errors
		err = handler(c, &tx)

		return err
	}
}
