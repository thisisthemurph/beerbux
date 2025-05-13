-- +goose Up
-- +goose StatementBegin
create table if not exists ledger (
    id uuid primary key default uuid_generate_v4(),
    transaction_id uuid not null references session_transactions(id) on delete no action,
    session_id uuid not null references sessions(id) on delete no action,
    user_id uuid not null references users(id) on delete no action,
    amount numeric(2,1) not null,
    created_at timestamp with time zone not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ledger;
-- +goose StatementEnd
