-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id uuid primary key default uuid_generate_v4(),
    username text not null unique,
    email text not null unique,
    update_email text,
    email_update_requested_at timestamp with time zone,
    email_update_otp text,
    email_last_updated_at timestamp with time zone,
    name text not null,
    hashed_password text not null,
    update_hashed_password text,
    password_update_requested_at timestamp with time zone,
    password_update_otp text,
    password_last_updated_at timestamp with time zone,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create index idx_users_username on users (username);

create trigger users_update_updated_at
    before update on users
    for each row
execute function fn_update_updated_at_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
