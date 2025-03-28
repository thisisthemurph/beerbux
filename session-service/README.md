# session-service

**Create a new session**

- Verify the user can create a session
- Create the session in the sessions table
- Add the user to the session_members table
- Send a session created event

## Manual testing

**Get a single session**

To get a single session, provide the ID of the session to be retrieved.

```shell
grpcurl -plaintext -d '{"session_id": "1c0327eb-b934-46be-a882-56195fab04d9"}' localhost:50052 session.service.Session.GetSession
````

**List all sessions for a user**

To list all sessions for a user, provide the ID of the user for which the sessions are to be listed.

```shell
grpcurl -plaintext -d '{"user_id": "10473635-01d4-4e2a-b809-8fce66031ace"}' localhost:50052 session.service.Session.ListSessionsForUser
```

```shell
grpcurl -plaintext -d '{"user_id": "10473635-01d4-4e2a-b809-8fce66031ace", "page_size": 1, "page_token": "4a861d2a-5848-4b33-b7b2-201c6cf74f4f"}' localhost:50052 session.service.Session.ListSessionsForUser
```

**Create a new session**

A session can be created by providing the ID of the creating user and the name of the session to be created.

```shell
grpcurl -plaintext -d '{"user_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde", "name": "My Test Session 2"}' localhost:50052 session.service.Session.CreateSession
```

**Add a user to a session**

A user can be added to a session by providing the ID of the user to be added and the ID of the session to which the user is to be added.

```shell
grpcurl -plaintext -d '{"user_id": "6cd0703c-1e23-43c6-96c2-af043e6ad4bf", "session_id": "5c0327eb-b934-46be-a882-56195fab04d9"}' localhost:50052 session.service.Session.AddMemberToSession
```

**Updating members of a session**

The session-service keeps its own copy of a sessions members. It records only those users who are members of a session, it does not keep a copy of all users.
Eventual consistency is ensured by listening to the `user.updated` event and updating the session members table.

Following is an example of an event for updating a member of a session:

```json
{
  "user_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde",
  "updated_fields": {
    "name": "New Name",
    "username": "new.username", 
    "updated_at": "2025-03-10T16:40:13Z"
  }
}
```