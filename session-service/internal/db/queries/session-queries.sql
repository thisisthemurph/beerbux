-- name: GetSession :one
select s.id, s.name, s.is_active, s.created_at, s.updated_at, cast(coalesce(sum(l.amount), 0) as real) as total
from sessions s
left join transactions t on s.id = t.session_id
left join transaction_lines l on t.id = l.transaction_id
where s.id = ?
group by s.id, s.name, s.is_active, s.created_at, s.updated_at;

-- name: ListSessionsForUser :many
with paged_sessions AS (
    select s.*, cast(coalesce(sum(l.amount), 0) as real) as total_amount
    from sessions s
    join session_members sm_target
        on s.id = sm_target.session_id
    left join transactions t
        on s.id = t.session_id
    left join transaction_lines l
        on t.id = l.transaction_id
    where sm_target.member_id = :member_id
        and (cast(coalesce(:page_token, '') as text) = '' or :page_token < s.id)
    group by s.id, s.name, s.is_active, s.created_at, s.updated_at
    order by s.updated_at desc, s.id desc
    limit case when :page_size = 0 then -1 else :page_size end
)
select
    ps.*,
    m.id as member_id,
    m.name as member_name,
    m.username as member_username
from paged_sessions ps
join session_members sm on ps.id = sm.session_id
join members m on sm.member_id = m.id
order by ps.updated_at desc, ps.id desc;

-- name: CreateSession :one
insert into sessions (id, name)
values (?, ?)
returning *;

-- name: GetMember :one
select * from members where id = ? limit 1;

-- name: ListMembers :many
select m.*
from members m
join session_members sm on m.id = sm.member_id
where sm.session_id = ?;

-- name: UpsertMember :exec
insert into members (id, name, username)
values (?, ?, ?)
on conflict(id) do update
set name = excluded.name,
    username = excluded.username,
    updated_at = current_timestamp;

-- name: UpdateMember :exec
update members
set name = ?,
    username = ?,
    updated_at = current_timestamp
where id = ?;

-- name: AddSessionMember :exec
insert into session_members (session_id, member_id, is_owner)
values (?, ?, ?)
on conflict do nothing;

-- name: UpdateSession :one
update sessions
set name = ?,
    updated_at = current_timestamp
where id = ?
returning *;

-- name: AddTransaction :one
insert into transactions (id, session_id, member_id)
values (?, ?, ?)
on conflict do nothing
returning *;

-- name: AddTransactionLine :one
insert into transaction_lines (transaction_id, member_id, amount)
values (?, ?, ?)
on conflict (transaction_id, member_id)
    do update set amount = excluded.amount
returning *;