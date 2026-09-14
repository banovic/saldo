package main

import (
	"context"
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

func applyMigrations(ctx context.Context, db *pgx.Conn, dir string, migrationFiles []string) error {
	const createMigrationsTableSQL = `
		CREATE TABLE IF NOT EXISTS migrations (
			filename text primary key,
			applied_at timestamptz not null default now()
		);
	`

	if _, err := db.Exec(ctx, createMigrationsTableSQL); err != nil {
		return err
	}

	applied, err := appliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	for _, file := range migrationFiles {
		if applied[file] {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return err
		}

		tx, err := db.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)

		if _, err = tx.Exec(ctx, string(content)); err != nil {
			return err
		}

		recordMigration := `
			INSERT INTO migrations (filename) VALUES ($1);
		`
		if _, err = tx.Exec(ctx, recordMigration, file); err != nil {
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func run(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var filesToProcess []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filesToProcess = append(filesToProcess, file.Name())
		}
	}

	if len(filesToProcess) == 0 {
		return nil
	}

	ctx := context.Background()

	dbUrl := os.Getenv("SALDO_DATABASE_URL")
	if dbUrl == "" {
		return fmt.Errorf("missing env var SALDO_DATABASE_URL")
	}

	db, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return err
	}

	defer db.Close(ctx)

	return applyMigrations(ctx, db, dir, filesToProcess)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dirname>\n", os.Args[0])
		os.Exit(2)
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Successfully applied migrations\n")
	os.Exit(0)
}
