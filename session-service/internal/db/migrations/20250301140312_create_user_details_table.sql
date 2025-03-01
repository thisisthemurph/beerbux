-- +goose Up
-- +goose StatementBegin
-- A shadow copy of the users table in the users.db
create table if not exists user_details (
    id text primary key not null,
    username text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_details;
-- +goose StatementEnd
