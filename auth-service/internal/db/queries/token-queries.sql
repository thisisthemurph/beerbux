-- name: RegisterRefreshToken :exec
insert into refresh_tokens (user_id, hashed_token, expires_at)
values (?, ?, ?);

-- name: GetRefreshTokensByUserID :many
select *
from refresh_tokens
where user_id = ?
and revoked = false
and expires_at > current_timestamp;

-- name: DeleteRefreshToken :exec
delete from refresh_tokens where id = ?;
