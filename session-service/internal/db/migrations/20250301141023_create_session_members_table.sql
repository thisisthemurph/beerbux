-- +goose Up
-- +goose StatementBegin
create table if not exists session_members (
    session_id text not null,
    user_id text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    primary key (session_id, user_id),
    foreign key (user_id) references user_details (id) on delete cascade,
    foreign key (session_id) references sessions (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists session_members;
-- +goose StatementEnd
