create table if not exists ledgers (
    ledger_id uuid primary key,
    name text not null,
    functional_currency text not null
);