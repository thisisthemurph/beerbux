-- +goose Up
-- +goose StatementBegin
create table if not exists transaction_lines (
    transaction_id text not null,
    member_id text not null,
    amount real not null default 0 check (amount >= 0),

    primary key (transaction_id, member_id),
    foreign key (transaction_id) references transactions (id) on delete cascade,
    foreign key (member_id) references members (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists transaction_lines;
-- +goose StatementEnd
