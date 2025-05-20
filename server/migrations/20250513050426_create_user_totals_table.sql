-- +goose Up
-- +goose StatementBegin
create table if not exists user_totals (
    user_id uuid references users(id) on delete cascade,
    credit numeric(5,1) not null default 0,
    debit numeric(5,1) not null default 0,
    primary key (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_totals;
-- +goose StatementEnd
