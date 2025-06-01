-- name: GetFriends :many
-- GetFriends returns the members that the provided member has sessions in common with.
-- The results include any sessions in which either of the members are deleted.
select u.id, u.name, u.username, count(distinct sm.session_id) as shared_session_count
from session_members sm
    join session_members sm2 on sm.session_id = sm2.session_id
    join users u on u.id = sm.member_id
where sm2.member_id = $1
    and sm.member_id != $1
group by u.id, u.username
order by shared_session_count desc, u.username;

-- name: MembersAreFriends :one
-- MembersAreFriends returns a boolean indicating if the provided members have any sessions in common.
-- This includes any sessions in which either of the members have been deleted from the session.
with joint_sessions as (
    select sm.session_id
    from session_members sm
    where sm.member_id = sqlc.arg(member_id)::uuid
       or sm.member_id = sqlc.arg(other_member_id)::uuid
    group by sm.session_id
    having count(distinct sm.member_id) = 2
)
select exists(select 1 from joint_sessions) as members_are_friends;

-- name: GetJointSessions :many
-- GetJointSessions returns details of the sessions for which both provided users are members.
-- If a user is a deleted member of a session, this session will not be returned.
with joint_sessions as (
    select sm.session_id
    from session_members sm
    where sm.is_deleted = false
        and (sm.member_id = sqlc.arg(member_id)::uuid
             or sm.member_id = sqlc.arg(other_member_id)::uuid)
    group by sm.session_id
    having count(distinct sm.member_id) = 2)
select s.*
from sessions s
join joint_sessions js on s.id = js.session_id;
