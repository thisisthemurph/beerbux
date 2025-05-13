-- +goose Up
-- +goose StatementBegin
create table if not exists session_transaction_lines (
    transaction_id uuid not null references session_transactions(id) on delete cascade,
    member_id uuid not null references users(id) on delete no action,
    amount numeric(2,1) not null default 0 check ( amount >= 0 ),
    primary key (transaction_id, member_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists session_transaction_lines;
-- +goose StatementEnd
