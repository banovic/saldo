package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
)

func appliedMigrations(ctx context.Context, db *pgx.Conn) (map[string]bool, error) {
	const query = `SELECT filename FROM migrations`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	filenames := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		filenames[filename] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return filenames, nil
}

func applyMigrations(ctx context.Context, db *pgx.Conn, dir string, migrationFiles []string) ([]string, error) {
	const createMigrationsTableSQL = `
		CREATE TABLE IF NOT EXISTS migrations (
			filename text primary key,
			applied_at timestamptz not null default now()
		);
	`

	const recordMigrationSQL = `
		INSERT INTO migrations (filename) VALUES ($1);
	`

	if _, err := db.Exec(ctx, createMigrationsTableSQL); err != nil {
		return nil, err
	}

	applied, err := appliedMigrations(ctx, db)
	if err != nil {
		return nil, err
	}

	var success []string

	for _, file := range migrationFiles {
		if applied[file] {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return success, fmt.Errorf("apply %s: %w", file, err)
		}

		err = pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, string(content)); err != nil {
				return err
			}

			_, err := tx.Exec(ctx, recordMigrationSQL, file)

			return err
		})
		if err != nil {
			return success, fmt.Errorf("apply %s: %w", file, err)
		}

		success = append(success, file)
	}
	return success, nil
}

func run(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filesToProcess []string
	for _, file := range entries {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filesToProcess = append(filesToProcess, file.Name())
		}
	}

	if len(filesToProcess) == 0 {
		return nil, nil
	}

	ctx := context.Background()

	dbURL := os.Getenv("SALDO_DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("missing env var SALDO_DATABASE_URL")
	}

	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	defer db.Close(ctx)

	return applyMigrations(ctx, db, dir, filesToProcess)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dirname>\n", os.Args[0])
		os.Exit(2)
	}

	success, err := run(os.Args[1])
	// Print executed migrations, as some could have failed.
	for _, file := range success {
		fmt.Printf("Applied migration: %s\n", file)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if len(success) == 0 {
		fmt.Printf("No migrations applied\n")
	}
}
