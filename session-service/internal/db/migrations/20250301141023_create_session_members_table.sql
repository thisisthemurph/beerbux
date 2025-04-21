-- +goose Up
-- +goose StatementBegin
create table if not exists session_members (
    session_id text not null,
    member_id text not null,
    is_owner boolean not null default false,
    is_admin boolean not null default false,
    is_deleted boolean not null default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    primary key (session_id, member_id),
    foreign key (member_id) references members (id) on delete cascade,
    foreign key (session_id) references sessions (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists session_members;
-- +goose StatementEnd
