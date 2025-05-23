-- +goose Up
-- +goose StatementBegin
create table if not exists session_members (
    session_id uuid not null references sessions(id) on delete cascade,
    member_id uuid not null references users(id) on delete no action,
    is_admin bool not null default false,
    is_deleted bool not null default false,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    primary key (session_id, member_id)
);

create trigger session_members_update_updated_at
    before update on session_members
    for each row
execute function fn_update_updated_at_timestamp();

create or replace function fn_session_members_trigger_sessions_mark_updated()
returns trigger as $$
begin
    perform fn_sessions_mark_updated(new.session_id);
    return new;
end;
$$ language plpgsql;

create trigger session_members_mark_sessions_updated
after insert or update on session_members
for each row
execute function fn_session_members_trigger_sessions_mark_updated();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists session_members;
-- +goose StatementEnd
