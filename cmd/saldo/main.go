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

	request := app.CreateLedgerRequest{Name: "Test6", FunctionalCurrency: "RSD", ReportingTimeZone: "Europe/Belgrade"}
	resp, err := service.CreateLedger(ctx, request)
	if err != nil {
		// TODO err needs to be mapped to response codes
		return err
	}
	fmt.Printf("%v\n", resp)

	req := app.GetLedgerRequest{LedgerID: "01a0e2e7-c0bc-75d4-a4fc-786ac7410be7"}
	resp2, err := service.GetLedger(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", resp2)

	req3 := app.ListLedgersRequest{}
	resp3, err := service.ListLedgers(ctx, req3)
	if err != nil {
		return err
	}
	for _, l := range resp3.Ledgers {
		fmt.Printf("%v\n", l)
	}
	return nil
}
