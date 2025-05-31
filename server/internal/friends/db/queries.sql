-- name: GetFriends :many
select u.id, u.name, u.username, count(distinct sm.session_id) as shared_session_count
from session_members sm
    join session_members sm2 on sm.session_id = sm2.session_id
    join users u on u.id = sm.member_id
where sm2.member_id = $1
    and sm.member_id != $1
group by u.id, u.username
order by shared_session_count desc, u.username;
