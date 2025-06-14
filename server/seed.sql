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
    ('connor', 'connor@example.com', 'CHC', '$2a$10$C7tUFbV.x7ZUzhjZLeDSBOhjSuXZhgoPP4OsnjjH33eR1nlgFou5.')
;

insert into sessions (name, creator_id)
values
    ('Ale House Boys', get_user_id('mike')),
    ('Empty Session', get_user_id('mike')),
    ('2 member session', get_user_id('connor')),
    ('Christmas Session', get_user_id('julian'))
;

insert into session_members (session_id, member_id, is_admin)
values
    (get_session_id('Ale House Boys'), get_user_id('mike'), true),
    (get_session_id('Ale House Boys'), get_user_id('julian'), false),
    (get_session_id('Ale House Boys'), get_user_id('andrew.longname'), false),
    (get_session_id('Ale House Boys'), get_user_id('connor'), false),
    (get_session_id('Empty Session'), get_user_id('mike'), true),
    (get_session_id('Christmas Session'), get_user_id('julian'), true),
    (get_session_id('Christmas Session'), get_user_id('mike'), true),
    (get_session_id('Christmas Session'), get_user_id('connor'), false),
    (get_session_id('2 member session'), get_user_id('connor'), true),
    (get_session_id('2 member session'), get_user_id('mike'), false)
;

-- Add a function to get transaction ID by session name and member (for consistency, optional)
create or replace function get_transaction_id(p_session_name text, p_username text)
    returns uuid
    language sql
as
$$
select st.id
from session_transactions st
         join sessions s on s.id = st.session_id
         join users u on u.id = st.member_id
where s.name = p_session_name and u.username = p_username
limit 1;
$$;

-- Insert into session_transactions
insert into session_transactions (session_id, member_id)
values
    (get_session_id('Ale House Boys'), get_user_id('mike')),
    (get_session_id('Ale House Boys'), get_user_id('julian')),
    (get_session_id('Christmas Session'), get_user_id('julian'))
returning id, session_id, member_id;

-- Assume the last inserted IDs represent:
-- tx1 (Mike pays 10.0): Mike paid for Julian and Connor
-- tx2 (Julian pays 6.0): Julian paid for Mike
-- tx3 (Julian pays 9.0): Julian paid for Mike and Connor

-- For clarity, you may want to capture the IDs with a script if using psql, but we'll use subqueries directly here

-- Insert session_transaction_lines for tx1
insert into session_transaction_lines (transaction_id, member_id, amount)
values
    ((select id from session_transactions where member_id = get_user_id('mike') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('mike'), 0.0),
    ((select id from session_transactions where member_id = get_user_id('mike') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('julian'), 5.0),
    ((select id from session_transactions where member_id = get_user_id('mike') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('connor'), 5.0);

-- Ledger entries for tx1
insert into ledger (transaction_id, user_id, amount)
values
    ((select id from session_transactions where member_id = get_user_id('mike') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('mike'), 10.0),
    ((select id from session_transactions where member_id = get_user_id('mike') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('julian'), -5.0),
    ((select id from session_transactions where member_id = get_user_id('mike') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('connor'), -5.0);

-- Insert session_transaction_lines for tx2
insert into session_transaction_lines (transaction_id, member_id, amount)
values
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('julian'), 0.0),
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('mike'), 6.0);

-- Ledger entries for tx2
insert into ledger (transaction_id, user_id, amount)
values
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('julian'), 6.0),
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Ale House Boys') limit 1), get_user_id('mike'), -6.0);

-- Insert session_transaction_lines for tx3 (Christmas Session)
insert into session_transaction_lines (transaction_id, member_id, amount)
values
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Christmas Session') limit 1), get_user_id('julian'), 0.0),
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Christmas Session') limit 1), get_user_id('mike'), 4.5),
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Christmas Session') limit 1), get_user_id('connor'), 4.5);

-- Ledger entries for tx3
insert into ledger (transaction_id, user_id, amount)
values
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Christmas Session') limit 1), get_user_id('julian'), 9.0),
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Christmas Session') limit 1), get_user_id('mike'), -4.5),
    ((select id from session_transactions where member_id = get_user_id('julian') and session_id = get_session_id('Christmas Session') limit 1), get_user_id('connor'), -4.5);

-- Cleanup helper functions

drop function if exists get_user_id;
drop function if exists get_session_id;
drop function if exists get_transaction_id;
