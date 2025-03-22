delete from refresh_tokens;
delete from users;


insert into users (id, username, hashed_password)
values
    ('10473635-01d4-4e2a-b809-8fce66031ace', 'mike', '$2a$10$/pYlYVC.cGJTfklMmrYICONMjiPoIKUZlW0UsyrfEOYVRsQR1tiuu'), -- password
    ('20473635-01d4-4e2a-b809-8fce66031ace', 'julian', '$2a$10$/pYlYVC.cGJTfklMmrYICONMjiPoIKUZlW0UsyrfEOYVRsQR1tiuu'), -- password
    ('30473635-01d4-4e2a-b809-8fce66031ace', 'andy', '$2a$10$/pYlYVC.cGJTfklMmrYICONMjiPoIKUZlW0UsyrfEOYVRsQR1tiuu'), -- password
    ('40473635-01d4-4e2a-b809-8fce66031ace', 'connor', '$2a$10$/pYlYVC.cGJTfklMmrYICONMjiPoIKUZlW0UsyrfEOYVRsQR1tiuu'); -- password
