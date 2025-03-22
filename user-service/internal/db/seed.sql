delete from user_ledger;
delete from users;

insert into users (id, name, username, bio, balance)
values
    ('10473635-01d4-4e2a-b809-8fce66031ace', 'Mike Murphy', 'mike', 'I am a software engineer and I like beer!', 26),
    ('20473635-01d4-4e2a-b809-8fce66031ace', 'Julian', 'julian', '', 26),
    ('30473635-01d4-4e2a-b809-8fce66031ace', 'Andy P', 'andy', null, 26),
    ('40473635-01d4-4e2a-b809-8fce66031ace', 'CHC', 'connor', null, 26);
