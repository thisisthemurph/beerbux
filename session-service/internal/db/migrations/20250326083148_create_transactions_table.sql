-- +goose Up
-- +goose StatementBegin
create table if not exists transactions (
    id text primary key,
    session_id text not null,
    member_id text not null,

    foreign key (session_id) references sessions (id) on delete cascade,
    foreign key (member_id) references members (id) on delete cascade
);

create index idx_transactions_session_id on transactions (session_id);
create index idx_transactions_member_id on transactions (member_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists transactions;
-- +goose StatementEnd
