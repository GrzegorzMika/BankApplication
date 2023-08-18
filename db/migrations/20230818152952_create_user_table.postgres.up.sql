create table if not exists users
(
    username            text primary key,
    hashed_password     text      not null,
    full_name           text      not null,
    email               text      not null unique,
    password_changed_at timestamp not null default '0001-01-01 00:00:00Z',
    created_at          timestamp not null default now()
);

alter table accounts
    add constraint owner_fk foreign key (owner) references users (username);

create unique index idx_owner_currency on accounts (owner, currency);