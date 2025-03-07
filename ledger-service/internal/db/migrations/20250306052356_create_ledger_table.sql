-- +goose Up
-- +goose StatementBegin
create table if not exists ledger (
    id text primary key,
    transaction_id text not null,
    session_id text not null,
    user_id text not null,
    amount real not null,
    created_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ledger;
-- +goose StatementEnd
