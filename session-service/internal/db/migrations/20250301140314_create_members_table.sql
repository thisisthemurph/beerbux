-- +goose Up
-- +goose StatementBegin
create table if not exists members (
    id text primary key not null,
    name text not null,
    username text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists members;
-- +goose StatementEnd
