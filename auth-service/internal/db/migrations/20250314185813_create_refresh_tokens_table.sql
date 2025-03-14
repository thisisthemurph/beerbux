-- +goose Up
-- +goose StatementBegin
create table if not exists refresh_tokens (
    id integer primary key,
    user_id text not null,
    hashed_token text not null unique,
    expires_at timestamp not null,
    revoked boolean not null default false,
    created_at timestamp not null default current_timestamp,
    foreign key (user_id) references users (id)
);

create index idx_refresh_tokens_user_revoked_expires_at on refresh_tokens (user_id, revoked, expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists refresh_tokens;
-- +goose StatementEnd
