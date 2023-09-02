create table if not exists sessions
(
    id uuid PRIMARY KEY,
    username            text not null,
    refresh_token     text      not null,
    user_agent           text      not null,
    client_ip              text      not null,
    is_blocked boolean not null,
    expires_at timestamp not null default '0001-01-01 00:00:00Z',
    created_at          timestamp not null default now()
);

alter table sessions
    add constraint sessions_fk foreign key (username) references users (username);