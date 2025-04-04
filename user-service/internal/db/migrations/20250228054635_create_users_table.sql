-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id text primary key not null,
    name text not null,
    username text not null unique,
    bio text,
    credit real not null default 0,
    debit real not null default 0,
    net real not null default 0,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
