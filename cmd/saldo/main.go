package main

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	errUsage     = errors.New("usage error")
	errFlagParse = errors.New("flag parse")
)

type command func(ctx context.Context, s *app.Service) (any, error)

func exitCode(err error) int {
	switch {
	case errors.Is(err, errUsage), errors.Is(err, errFlagParse):
		return 1
	case errors.Is(err, app.ErrInvalidInput):
		return 2
	case errors.Is(err, app.ErrNotFound):
		return 3
	case errors.Is(err, app.ErrInternal):
		return 4
	default:
		return 100
	}
}

func parseCommand(args []string) (command, error) {
	switch args[0] {
	case "create_ledger":
		var req app.CreateLedgerRequest
		fs := flag.NewFlagSet("create_ledger", flag.ContinueOnError)
		fs.StringVar(&req.Name, "name", "", "Ledger name")
		fs.StringVar(&req.FunctionalCurrency, "functional_currency", "", "Functional currency")
		fs.StringVar(&req.ReportingTimeZone, "reporting_time_zone", "", "Reporting time zone IANA format")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, fmt.Errorf("%w: %w", errFlagParse, err)
		}
		return func(ctx context.Context, s *app.Service) (any, error) {
			return s.CreateLedger(ctx, req)
		}, nil
	case "get_ledger":
		var req app.GetLedgerRequest
		fs := flag.NewFlagSet("get_ledger", flag.ContinueOnError)
		fs.StringVar(&req.LedgerID, "ledger_id", "", "Ledger id")
		if err := fs.Parse(args[1:]); err != nil {
			return nil, fmt.Errorf("%w: %w", errFlagParse, err)
		}
		return func(ctx context.Context, s *app.Service) (any, error) {
			return s.GetLedger(ctx, req)
		}, nil
	case "list_ledgers":
		var req app.ListLedgersRequest
		fs := flag.NewFlagSet("list_ledgers", flag.ContinueOnError)
		if err := fs.Parse(args[1:]); err != nil {
			return nil, fmt.Errorf("%w: %w", errFlagParse, err)
		}
		return func(ctx context.Context, s *app.Service) (any, error) {
			return s.ListLedgers(ctx, req)
		}, nil
	default:
		return nil, fmt.Errorf("%w: unrecognized command: %q", errUsage, args[0])
	}
}

const usage = `Usage: SALDO_DATABASE_URL=<url> saldo <command> [flags]

Commands:
  create_ledger  Create a new ledger.
  get_ledger     Show one ledger.
  list_ledgers   List all ledgers.

Run 'saldo <command> -h' to see the command's flags.
`

func printUsage() {
	fmt.Fprint(os.Stderr, usage)
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("%w: command required", errUsage)
	}

	command, err := parseCommand(os.Args[1:])
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Build service.

	dbURL := os.Getenv("SALDO_DATABASE_URL")
	if dbURL == "" {
		return errors.New("missing env var SALDO_DATABASE_URL")
	}

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}

	defer dbpool.Close()

	service := app.NewService(postgres.NewUnitOfWork(dbpool), time.Now, uuid.NewV7)

	resp, err := command(ctx, service)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	switch {
	case err == nil:
		// Do nothing, exit 0
	case errors.Is(err, flag.ErrHelp):
		// Help was asked for (-h / --help), and FlagSet printed usage.
		// Do nothing, exit 0
	case errors.Is(err, errUsage):
		// Usage error - print error and then usage.
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(exitCode(err))
	case errors.Is(err, errFlagParse):
		// FlagSet printed usage, exit with error (this is not usage in the non-error sense).
		os.Exit(exitCode(err))
	default:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode(err))
	}
}
