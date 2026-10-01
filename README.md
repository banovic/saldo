# saldo

A double-entry accounting ledger in Go. Data is stored in PostgreSQL.
saldo records financial events as balanced journal entries and never rewrites history.
Mistakes are corrected by new entries reversing wrong ones.
Account's balance is always calculated from its postings and is never stored nor updated.

Status: in progress, domain model and testing in place.


## Invariants

These invariants are enforced in the domain layer, not left to callers or the database.

**Balanced entries** - A journal entry has at least two postings, and their functional amounts sum to exactly zero in the ledger's functional currency.

**Append-only history** - Journal entries are never updated or deleted. A correction is a new entry that references the one it reverses; an entry cannot reverse itself.

**Exactly-once recording** - Every journal entry carries a client-supplied idempotency key, unique per ledger. Uniqueness will be enforced by database check.

**Safe money arithmetic** - Money is represented as int64 minor units plus an ISO 4217 currency; the currency's exponent defines the minor unit. Arithmetic is overflow-checked and summing mixed currencies is an error.

**Derived balances** - Accounts store no balance; a balance is the sum of the account's postings. Accounts are typed by the accounting equation: asset, liability, equity, revenue, expense.


## Architecture

Architecture is ports and adapters. Full layering rules are in doc.go.
Layers are work in progress

| Layer | Responsibility
| ----- | --------------
|domain	|Entities, value objects, invariants, repository interfaces. Standard library only.
|app |One use case per file; request/response DTOs. All database access goes through a UnitOfWork, which is the only transaction boundary. Clock and ID generation are injected.
|infrastructure/postgres |UnitOfWork owns the transaction; repositories never begin or commit. Postgres driver errors are translated into app errors here.
|cmd/saldo, cmd/saldod |CLI and HTTP (not yet implemented) runners: config, manual constructor-injection wiring, and mapping app error codes to exit codes or HTTP status codes.
|cmd/migrate| Minimal forward-only migrator: applies schema/NNNN_*.sql in order, one transaction per file.

IDs are UUIDv7. Dependencies are kept deliberately minimal (pgx for Postgres); tests use the standard library only, table-driven.



## Run it

Start up Postgres server:

```
docker run -d --name saldo-pg -p 5432:5432 -e POSTGRES_PASSWORD=saldo postgres:18
```

apply database migrations:

```
SALDO_DATABASE_URL=postgres://postgres:saldo@localhost:5432/postgres go run ./cmd/migrate schema
```

then start up application:

```
SALDO_DATABASE_URL=postgres://postgres:saldo@localhost:5432/postgres go run ./cmd/saldo
```

Open a Postgres client in the container:

```
docker exec -it saldo-pg psql -U postgres
```

## Glossary
'to book' to record in the books.
