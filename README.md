# saldo

A double-entry accounting ledger in Go.

## Run it

Start up Postgres server:

```
docker run -d --name saldo-pg -p 5432:5432 -e POSTGRES_PASSWORD=saldo postgres:18
```

then start up application:

```
SALDO_DATABASE_URL=postgres://postgres:saldo@localhost:5432/postgres go run ./cmd/saldo
```

## Glossary
'to book' to record in the books.
