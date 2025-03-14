-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id text not null primary key,
    username text not null,
    hashed_password text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index idx_users_username on users (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
