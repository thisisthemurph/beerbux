-- name: Exists :one
select exists(select 1 from sessions where id = $1);

-- name: Get :one
select s.id, s.name, s.is_active, s.created_at, s.updated_at, coalesce(sum(l.amount), 0)::float8 as total
from sessions s
left join session_transactions t on s.id = t.session_id
left join session_transaction_lines l on t.id = l.transaction_id
where s.id = $1
group by s.id, s.name, s.is_active, s.created_at, s.updated_at;

-- name: ListSessionsForUser :many
with paged_sessions as (
    select
        s.id,
        s.name,
        s.is_active,
        s.created_at,
        s.updated_at,
        coalesce(sum(tl.amount), 0)::float8 as total_amount
    from sessions s
             join session_members sm_target on s.id = sm_target.session_id
             left join session_transactions t on s.id = t.session_id
             left join session_transaction_lines tl on t.id = tl.transaction_id
    where sm_target.member_id = $1
      and sm_target.is_deleted = false
    group by s.id, s.name, s.is_active, s.created_at, s.updated_at
    order by s.updated_at desc, s.id desc
    limit nullif($2, 0)
)
select
    st.*,
    m.id as member_id,
    m.name as member_name,
    m.username as member_username
from paged_sessions st
join session_members sm on st.id = sm.session_id
join users m on sm.member_id = m.id
order by st.updated_at desc, st.id desc;

