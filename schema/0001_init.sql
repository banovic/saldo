create table ledgers (
    ledger_id uuid primary key,
    name text not null,
    functional_currency text not null,
    reporting_time_zone text not null,
    created_at timestamptz not null
);

create unique index ledgers_name_unique ON ledgers (lower(name));

create type account_type as enum (
    'asset', 'liability', 'equity', 'revenue', 'expense'
);

create table accounts (
    account_id uuid primary key,
    ledger_id uuid not null references ledgers(ledger_id) on delete restrict,
    type account_type not null,
    code text not null check(length(code) <= 200),
    name text not null check(length(name) <= 200),
    created_at timestamptz not null
);

create table journal_entries (
    journal_entry_id uuid primary key,
    ledger_id uuid not null references ledgers(ledger_id) on delete restrict,
    idempotency_key text not null check(length(idempotency_key) <= 200),
    description text not null check(length(description) <= 1000),
    reverses_id uuid references journal_entries(journal_entry_id),
    occurred_at timestamptz not null,
    created_at timestamptz not null,
    posted_on date not null
);

create unique index journal_entries_idempotency_key_unique on journal_entries (ledger_id, idempotency_key);

create table exchange_rates (
    exchange_rate_id uuid primary key,
    base_currency text not null check(length(base_currency) = 3),
    quote_currency text not null check(length(quote_currency) = 3),
    num bigint not null check(num > 0),
    den bigint not null check(den > 0),
    source text not null check(length(source) <= 200),
    fixed_on date not null,
    created_at timestamptz not null,
    check(base_currency <> quote_currency)
);

create table postings (
    posting_id uuid primary key,
    journal_entry_id uuid not null references journal_entries(journal_entry_id) on delete restrict,
    account_id uuid not null references accounts(account_id) on delete restrict,
    exchange_rate_id uuid references exchange_rates(exchange_rate_id) on delete restrict,
    transaction_minor_units bigint not null,
    transaction_currency text not null check(length(transaction_currency) = 3),
    functional_minor_units bigint not null,
    functional_currency text not null check(length(functional_currency) = 3)
);

create index postings_journal_entry_id_idx on postings(journal_entry_id);
create index postings_account_id_idx on postings(account_id);
