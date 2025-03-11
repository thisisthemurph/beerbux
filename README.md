# Beerbux

Beerbux is a distributed system allowing you to track beers in and beers out.

- Create a `session` with friends and track who is buying the beer, find out who owes who and how much.
- Keep track of beers owed in a session as well as your overall debt and credit.

## Running locally

Start the Kafka cluster in Docker:

```shell
docker compose up -d
```

**Databases**

Each of the services persists their data in their own databases.
No service has direct access to a database belonging to another service.
