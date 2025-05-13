-- name: Get :one
select u.id, u.username, u.name, u.created_at, u.updated_at,
       t.debit, t.credit
from users u
left join user_totals t on u.id = t.user_id
where u.id = $1
limit 1;

-- name: GetByUsername :one
select u.id, u.username, u.name, u.created_at, u.updated_at,
       t.debit, t.credit
from users u
left join user_totals t on u.id = t.user_id
where u.username = $1
limit 1;

-- name: Update :one
update users
set name = $1, username = $2
where id = $3
returning id, username, name, created_at, updated_at;

-- name: UpdateTotals :exec
insert into user_totals (user_id, credit, debit)
values ($1, $2, $3)
on conflict (user_id) do update
set credit = excluded.credit, debit = excluded.debit;

-- name: UserWithUsernameExists :one
select exists(select 1 from users where username = $1);
