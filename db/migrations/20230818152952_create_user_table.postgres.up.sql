create table if not exists users
(
    username            text primary key,
    hashed_password     text      not null,
    full_name           text      not null,
    email               text      not null unique,
    is_email_verified   boolean   not null default false,
    password_changed_at timestamp not null default '0001-01-01 00:00:00Z',
    created_at          timestamp not null default now()
);

alter table accounts
    add constraint owner_fk foreign key (owner) references users (username);

create unique index idx_owner_currency on accounts (owner, currency);

create table if not exists verify_emails
(
    id          serial primary key,
    username    varchar   not null references users (username),
    email       varchar   not null,
    secret_code varchar   not null,
    is_used     boolean   not null default false,
    created_at  timestamp not null default now(),
    expired_at  timestamp not null default now() + interval '15 minutes'
);