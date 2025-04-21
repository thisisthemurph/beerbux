-- +goose Up
-- +goose StatementBegin
create table if not exists session_history(
    id integer primary key,
    session_id text not null,
    member_id text not null,
    event_type text not null,
    event_data blob,
    created_at timestamp not null default current_timestamp,
    foreign key (session_id) references sessions(id),
    foreign key (member_id) references members(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists session_history;
-- +goose StatementEnd
