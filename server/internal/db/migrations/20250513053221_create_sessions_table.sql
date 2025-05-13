-- +goose Up
-- +goose StatementBegin
create table if not exists sessions (
  id uuid primary key default uuid_generate_v4(),
  name text not null,
  is_active bool not null default true,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now()
);

create trigger sessions_update_updated_at
    before update on sessions
    for each row
execute function fn_update_updated_at_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists sessions;
-- +goose StatementEnd
