-- +goose Up
-- +goose StatementBegin
create table if not exists session_transactions (
    id uuid primary key not null default uuid_generate_v4(),
    session_id uuid not null references sessions(id) on delete no action,
    member_id uuid not null references users(id) on delete no action,
    created_at timestamp with time zone not null default now()
);

create or replace function fn_session_transactions_trigger_sessions_mark_updated()
    returns trigger as $$
begin
    perform fn_sessions_mark_updated(new.session_id);
    return new;
end;
$$ language plpgsql;

create trigger session_members_mark_sessions_updated
after insert on session_transactions
for each row
execute function fn_session_transactions_trigger_sessions_mark_updated();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists session_transactions;
-- +goose StatementEnd
