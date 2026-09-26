package main

import (
	"context"
	"fmt"
	"os"
	"time"

	// time.LoadLocation needs tzdata on the host, and some minimal images might not have it.
	// This will embed tzdata in built binary, so no problems on such hosts.
	_ "time/tzdata"
	"uuid"

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

	service := app.NewService(postgres.NewUnitOfWork(dbpool), time.Now, uuid.NewV7)

	request := app.CreateLedgerRequest{Name: "Test", FunctionalCurrency: "RSD", ReportingTimeZone: "Europe/Belgrade"}
	resp, err := service.CreateLedger(ctx, request)
	if err != nil {
		// TODO err needs to be mapped to response codes
		return err
	}
	fmt.Printf("%v\n", resp)
	return nil
}
