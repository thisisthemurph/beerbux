-- +goose Up
-- +goose StatementBegin
create table if not exists sessions (
    id text primary key not null,
    name text not null,
    owner_id text not null,
    is_active boolean not null default true,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    foreign key (owner_id) references user_details (id) on delete set null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists sessions;
-- +goose StatementEnd
