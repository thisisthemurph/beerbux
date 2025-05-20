-- +goose Up
-- +goose StatementBegin
create or replace function fn_update_user_totals()
    returns trigger
    language plpgsql
as
$$
begin
    insert into user_totals (user_id, credit, debit)
    select
        user_id,
        coalesce(sum(case when amount > 0 then amount else 0 end), 0) as credit,
        coalesce(sum(case when amount < 0 then abs(amount) else 0 end), 0) as debit
    from new_ledger_rows
    group by user_id
    on conflict (user_id) do update
        set
            credit = user_totals.credit + excluded.credit,
            debit = user_totals.debit + excluded.debit;

    return null;
end;
$$;

create trigger tr_update_user_totals
    after insert on ledger
    referencing new table as new_ledger_rows
    for each statement
execute function fn_update_user_totals();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger if exists tr_update_user_totals on ledger;
drop function if exists fn_update_user_totals;
-- +goose StatementEnd
