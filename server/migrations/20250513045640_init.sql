-- +goose Up
-- +goose StatementBegin
create extension if not exists "uuid-ossp";

create or replace function fn_update_updated_at_timestamp()
    returns trigger as $$
begin
    new.updated_at = current_timestamp;
return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop extension "uuid-ossp";
drop function fn_update_updated_at_timestamp;
-- +goose StatementEnd
