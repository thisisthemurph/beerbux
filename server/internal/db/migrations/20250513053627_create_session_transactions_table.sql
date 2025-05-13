-- +goose Up
-- +goose StatementBegin
create table if not exists session_transactions (
    id uuid primary key not null default uuid_generate_v4(),
    session_id uuid not null references sessions(id) on delete no action,
    member_id uuid not null references users(id) on delete no action,
    created_at timestamp with time zone not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists session_transactions;
-- +goose StatementEnd
