package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/banovic/saldo/app"
	"github.com/banovic/saldo/infrastructure/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("SALDO_DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	service := app.NewService(postgres.NewUnitOfWork(dbpool), time.Now)

	request := app.CreateLedgerRequest{}
	resp, err := service.CreateLedger(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", resp)
}
