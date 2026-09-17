create table if not exists ledgers (
    ledger_id uuid primary key,
    name text not null,
    functional_currency text not null,
    reporting_time_zone text not null,
    created_at timestamptz not null
);