- Minimal dependencies. Ask before adding anything.
- No assertion libraries. Table-driven tests with stdlib testing.
- Be brief. Answer the question asked, then stop. No multi-section essays.

## Project

Double-entry accounting ledger in Go (`github.com/banovic/saldo`). Only dep: pgx/v5.
Early stage: domain model is real and tested; app/infrastructure/cmd are wiring skeletons.

## Commands

- `go test ./...` — only `domain` has tests.
- `go vet ./...`
- Run: see README (Postgres in docker, `SALDO_DATABASE_URL`). No schema/migrations yet.

## Architecture

Layering is specified in [doc.go](doc.go) — read it before structural changes.

- `domain` — entities, value types, repository interfaces. No imports outside stdlib.
- `app` — one use case per file on `*Service`, request/response DTOs. All DB access via
  `UnitOfWork.Execute(ctx, func(Repositories) error)`. Clock injected as `now func() time.Time`.
- `infrastructure/postgres` — `UnitOfWork` owns the tx; repositories hold `pgx.Tx`, never begin/commit.
  Translates driver errors to app errors.
- `cmd/saldo` (CLI), `cmd/saldod` (HTTP) — config, manual constructor-injection wiring, map app error codes to exit/HTTP codes.

Not built yet: `app/error.go`, `cmd/saldo/config.go` (empty), `cmd/saldod/config.go` (missing),
repository `Insert`s (return TODO), `CreateLedger` (stub), HTTP runner.

## Domain rules

- Money = `int64` minor units + `Currency`. Guard overflow (`ErrOverflow`); sums via `big.Int`.
- Currency = ISO 4217 code; known set lives in `currencyInfo` map (exponent drives minor units).
- `JournalEntry`: ≥2 postings, `FunctionalAmount`s sum to zero in ledger's functional currency.
  Immutable — corrections are new entries with `Reverses` set. `IdempotencyKey` for exactly-once.
- Posting carries `TransactionAmount` (original currency) and `FunctionalAmount` (booked) + `ExchangeRateID`.
- `ExchangeRate` is append-only, stored as `Num/Den` integers, with `Kind` and `Source`.
- Account has no stored balance; balance = sum of its postings' functional amounts.

## Conventions

- ID/enum types: `IsValid() bool`, with doc comment listing validity rules. Aggregates: `Validate(...) error`.
- Sentinel `Err*` vars wrapped with `fmt.Errorf("%w: ...")`; tests check with `errors.Is`.
- Tests: `testCases := []struct{name ...}` + `t.Run`, messages as `f(x) = got, want want`.
- Interface assertions: `var _ app.UnitOfWork = (*UnitOfWork)(nil)`.
