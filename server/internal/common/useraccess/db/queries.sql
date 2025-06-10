-- name: GetUserByID :one
select
    u.id, u.username, u.email, u.name, u.created_at, u.updated_at,
    coalesce(ut.debit, 0) as debit,
    coalesce(ut.credit, 0) as credit,
    coalesce(ucs.credit_score, 0) as credit_score
from users u
left join user_totals ut on u.id = ut.user_id
left join user_credit_score ucs on u.id = ucs.user_id
where u.id = $1
limit 1;

-- name: GetByUsername :one
select
    u.id, u.username, u.email, u.name, u.created_at, u.updated_at,
    coalesce(ut.debit, 0) as debit,
    coalesce(ut.credit, 0) as credit,
    coalesce(ucs.credit_score, 0) as credit_score
from users u
left join user_totals ut on u.id = ut.user_id
left join user_credit_score ucs on u.id = ucs.user_id
where u.username = $1
limit 1;

-- name: GetUserByEmail :one
select
    u.id, u.username, u.email, u.name, u.created_at, u.updated_at,
    coalesce(ut.debit, 0) as debit,
    coalesce(ut.credit, 0) as credit,
    coalesce(ucs.credit_score, 0) as credit_score
from users u
left join user_totals ut on u.id = ut.user_id
left join user_credit_score ucs on u.id = ucs.user_id
where u.email = $1
limit 1;

-- name: UserWithUsernameExists :one
select exists(select 1 from users where username = $1);

-- name: UserWithEmailExists :one
select exists(select 1 from users where email = $1);

-- name: GetUserCreditScore :one
select * from user_credit_score where user_id = $1 limit 1;

