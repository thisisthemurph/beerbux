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
