-- noinspection SqlWithoutWhereForFile
begin;
delete from ledger;
delete from session_history;
delete from session_transaction_lines;
delete from session_transactions;
delete from session_members;
delete from sessions;
delete from user_totals;
delete from users;
commit;

create or replace function get_user_id(p_username text)
returns uuid
language sql
as
$$
    select id from users where username = p_username limit 1;
$$;

create or replace function get_session_id(p_session_name text)
returns uuid
language sql
as
$$
    select id from sessions where name = p_session_name limit 1;
$$;

insert into users (username, email, name, hashed_password)
values
    ('mike', 'mike@example.com', 'Mike', '$2a$10$C7tUFbV.x7ZUzhjZLeDSBOhjSuXZhgoPP4OsnjjH33eR1nlgFou5.'),
    ('julian', 'julian@example.com', 'Julian', '$2a$10$C7tUFbV.x7ZUzhjZLeDSBOhjSuXZhgoPP4OsnjjH33eR1nlgFou5.'),
    ('andrew.longname', 'andy@example.com', 'Andrew Longname', '$2a$10$C7tUFbV.x7ZUzhjZLeDSBOhjSuXZhgoPP4OsnjjH33eR1nlgFou5.'),
    ('connor', 'connor@example.com', 'CHC', '$2a$10$C7tUFbV.x7ZUzhjZLeDSBOhjSuXZhgoPP4OsnjjH33eR1nlgFou5.');

insert into sessions (name, creator_id)
values
    ('Ale House Boys', get_user_id('mike')),
    ('Christmas Session', get_user_id('julian'));

insert into session_members (session_id, member_id, is_admin)
values
    (get_session_id('Ale House Boys'), get_user_id('mike'), true),
    (get_session_id('Ale House Boys'), get_user_id('julian'), false),
    (get_session_id('Ale House Boys'), get_user_id('andrew.longname'), false),
    (get_session_id('Ale House Boys'), get_user_id('connor'), false),
    (get_session_id('Christmas Session'), get_user_id('julian'), true),
    (get_session_id('Christmas Session'), get_user_id('mike'), true),
    (get_session_id('Christmas Session'), get_user_id('connor'), false);

drop function if exists get_user_id;
drop function if exists get_session_id;
