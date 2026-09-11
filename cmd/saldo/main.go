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
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, os.Getenv("SALDO_DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}

	defer dbpool.Close()

	service := app.NewService(postgres.NewUnitOfWork(dbpool), time.Now)

	request := app.CreateLedgerRequest{}
	resp, err := service.CreateLedger(ctx, request)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", resp)
	return nil
}
