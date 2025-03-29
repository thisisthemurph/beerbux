delete from transaction_lines;
delete from transactions;
delete from session_members;
delete from members;
delete from sessions;

insert into members (id, name, username)
values
    ('10473635-01d4-4e2a-b809-8fce66031ace', 'Mike Murphy', 'mike'),
    ('20473635-01d4-4e2a-b809-8fce66031ace', 'Julian', 'julian'),
    ('30473635-01d4-4e2a-b809-8fce66031ace', 'Andy P', 'andy'),
    ('40473635-01d4-4e2a-b809-8fce66031ace', 'CHC', 'connor');

insert into sessions (id, name, is_active)
values
    ('1c0327eb-b934-46be-a882-56195fab04d9', 'Ale House Boys', true),
    ('2c0327eb-b934-46be-a882-56195fab04d9', 'Christmas Sesh', false),
    ('3c0327eb-b934-46be-a882-56195fab04d9', 'Random session', true),
    ('4c0327eb-b934-46be-a882-56195fab04d9', 'Lem War Crawl', true);

insert into session_members (session_id, member_id, is_owner)
values
    ('1c0327eb-b934-46be-a882-56195fab04d9', '10473635-01d4-4e2a-b809-8fce66031ace', true),
    ('1c0327eb-b934-46be-a882-56195fab04d9', '20473635-01d4-4e2a-b809-8fce66031ace', false),
    ('1c0327eb-b934-46be-a882-56195fab04d9', '30473635-01d4-4e2a-b809-8fce66031ace', false),
    ('1c0327eb-b934-46be-a882-56195fab04d9', '40473635-01d4-4e2a-b809-8fce66031ace', false),

    ('2c0327eb-b934-46be-a882-56195fab04d9', '20473635-01d4-4e2a-b809-8fce66031ace', true),
    ('2c0327eb-b934-46be-a882-56195fab04d9', '30473635-01d4-4e2a-b809-8fce66031ace', false),
    ('2c0327eb-b934-46be-a882-56195fab04d9', '40473635-01d4-4e2a-b809-8fce66031ace', false),

    ('3c0327eb-b934-46be-a882-56195fab04d9', '10473635-01d4-4e2a-b809-8fce66031ace', false),
    ('3c0327eb-b934-46be-a882-56195fab04d9', '20473635-01d4-4e2a-b809-8fce66031ace', false),
    ('3c0327eb-b934-46be-a882-56195fab04d9', '30473635-01d4-4e2a-b809-8fce66031ace', true),

    ('4c0327eb-b934-46be-a882-56195fab04d9', '10473635-01d4-4e2a-b809-8fce66031ace', false),
    ('4c0327eb-b934-46be-a882-56195fab04d9', '20473635-01d4-4e2a-b809-8fce66031ace', true),
    ('4c0327eb-b934-46be-a882-56195fab04d9', '30473635-01d4-4e2a-b809-8fce66031ace', false),
    ('4c0327eb-b934-46be-a882-56195fab04d9', '40473635-01d4-4e2a-b809-8fce66031ace', false);

insert into transactions (id, session_id, member_id)
values
    ('f3d7ec98-034d-4c25-8b9e-023faa19fd37', '1c0327eb-b934-46be-a882-56195fab04d9', '10473635-01d4-4e2a-b809-8fce66031ace'),
    ('ed36a838-db7b-4de6-93f5-5bcde30a1a30', '1c0327eb-b934-46be-a882-56195fab04d9', '10473635-01d4-4e2a-b809-8fce66031ace'),
    ('1d732ace-5034-483c-95b1-879a75c297f5', '1c0327eb-b934-46be-a882-56195fab04d9', '20473635-01d4-4e2a-b809-8fce66031ace');

insert into transaction_lines (transaction_id, member_id, amount)
values
    ('f3d7ec98-034d-4c25-8b9e-023faa19fd37', '20473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('f3d7ec98-034d-4c25-8b9e-023faa19fd37', '30473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('f3d7ec98-034d-4c25-8b9e-023faa19fd37', '40473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('ed36a838-db7b-4de6-93f5-5bcde30a1a30', '20473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('ed36a838-db7b-4de6-93f5-5bcde30a1a30', '30473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('ed36a838-db7b-4de6-93f5-5bcde30a1a30', '40473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('1d732ace-5034-483c-95b1-879a75c297f5', '10473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('1d732ace-5034-483c-95b1-879a75c297f5', '30473635-01d4-4e2a-b809-8fce66031ace', 1),
    ('1d732ace-5034-483c-95b1-879a75c297f5', '40473635-01d4-4e2a-b809-8fce66031ace', 1);
