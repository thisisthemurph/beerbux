-- +goose Up
-- +goose StatementBegin
create table if not exists sessions (
    id uuid primary key default uuid_generate_v4(),
    name text not null,
    is_active bool not null default true,
    creator_id uuid not null references users(id) on delete no action,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create trigger sessions_update_updated_at
    before update on sessions
    for each row
execute function fn_update_updated_at_timestamp();

create or replace function fn_sessions_mark_updated(session_id uuid)
returns void as $$
begin
    update sessions
    set updated_at = now()
    where id = session_id;
end;
$$ language plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists sessions;
-- +goose StatementEnd
