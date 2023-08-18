CREATE TABLE IF NOT EXISTS
(
    id         bigserial primary key,
    owner      varchar not null,
    balance    float   not null,
    currency   char(3) not null,
    created_at timestamp with time zone default now()
);

CREATE INDEX IF NOT EXISTS owner_idx ON accounts (owner);

CREATE TABLE IF NOT EXISTS entries
(
    id         bigserial primary key,
    account_id bigint not null references accounts (id) on delete cascade,
    amount     float  not null,
    created_at timestamp with time zone default now()
);

CREATE INDEX IF NOT EXISTS account_idx ON entries (account_id);

CREATE TABLE IF NOT EXISTS transfers
(
    id              bigserial primary key,
    from_account_id bigint not null references accounts (id) on delete cascade,
    to_account_id   bigint not null references accounts (id) on delete cascade,
    amount          float  not null check (amount > 0),
    created_at      timestamp with time zone default now()
);

CREATE INDEX IF NOT EXISTS from_account_id_idx ON transfers (from_account_id);
CREATE INDEX IF NOT EXISTS to_account_id_idx ON transfers (to_account_id);
CREATE INDEX IF NOT EXISTS from_to_account_id_idx ON transfers (from_account_id, to_account_id);

