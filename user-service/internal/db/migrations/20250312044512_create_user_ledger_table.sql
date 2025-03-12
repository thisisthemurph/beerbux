-- +goose Up
-- +goose StatementBegin
create table if not exists user_ledger (
    id integer primary key,
    user_id text not null,
    participant_id text not null,
    amount real not null,
    type text not null,
    created_at timestamp not null default current_timestamp,
    foreign key (user_id) references users (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_ledger;
-- +goose StatementEnd
