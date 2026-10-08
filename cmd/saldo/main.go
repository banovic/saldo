package main

import (
	"context"
	"flag"
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

type command func(ctx context.Context, s *app.Service) (any, error)

func parseCommand(args []string) (command, error) {
	switch args[0] {
	case "create_ledger":
		var req app.CreateLedgerRequest
		fs := flag.NewFlagSet("create_ledger", flag.ContinueOnError)
		fs.StringVar(&req.Name, "name", "", "Ledger name")
		fs.StringVar(&req.FunctionalCurrency, "functional_currency", "", "Functional currency")
		fs.StringVar(&req.ReportingTimeZone, "reporting_time_zone", "", "Reporting time zone IANA format")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		return func(ctx context.Context, s *app.Service) (any, error) {
			return s.CreateLedger(ctx, req)
		}, nil
	case "get_ledger":
		var req app.GetLedgerRequest
		fs := flag.NewFlagSet("get_ledger", flag.ContinueOnError)
		fs.StringVar(&req.LedgerID, "ledger_id", "", "Ledger id")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, err
		}
		return func(ctx context.Context, s *app.Service) (any, error) {
			return s.GetLedger(ctx, req)
		}, nil
	case "list_ledgers":
		var req app.ListLedgersRequest
		return func(ctx context.Context, s *app.Service) (any, error) {
			return s.ListLedgers(ctx, req)
		}, nil
	default:
		return nil, fmt.Errorf("Unrecognized command: %s", args[0])
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Expected at least 2 arguments")
	}

	ctx := context.Background()

	// Build service.

	dbURL := os.Getenv("SALDO_DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("missing env var SALDO_DATABASE_URL")
	}

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}

	defer dbpool.Close()

	service := app.NewService(postgres.NewUnitOfWork(dbpool), time.Now, uuid.NewV7)

	command, err := parseCommand(os.Args[1:])
	if err != nil {
		return err
	}

	resp, err := command(ctx, service)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", resp)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
