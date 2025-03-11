# ledger-service

The ledger-service listens for `transaction.created` events, validates that ledger, 
and publishes events detailing the ledger with the newly created ledger entries.

## Manual testing

A `transaction.created` event can be published to Kafka.
This will be collected; entries will be created in the ledger table and `ledger.updated` events will be published.

```json
{
  "transaction_id": "c6037d00-88e2-4479-9468-fb79033fcd27",
  "creator_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde",
  "session_id": "5c0327eb-b934-46be-a882-56195fab04d9",
  "member_amounts": [
    {
      "user_id": "6cd0703c-1e23-43c6-96c2-af043e6ad4bf",
      "amount": 1
    }
  ]
}
```

The above event will create two ledger entries in the database, one for the creator and one for the member.
Two `ledger.updated` events will be published.

One for the creator of the transaction:

```json
{
	"id": "1800615e-fd8a-4600-8a27-7e056fe940be",
	"transaction_id": "c6037d00-88e2-4479-9468-fb79033fcd27",
	"session_id": "5c0327eb-b934-46be-a882-56195fab04d9",
	"user_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde",
	"amount": -1
}
```

And one for the member of the transaction:

```json
{
	"id": "e6140322-ed5a-497b-a394-8b511908d181",
	"transaction_id": "c6037d00-88e2-4479-9468-fb79033fcd27",
	"session_id": "5c0327eb-b934-46be-a882-56195fab04d9",
	"user_id": "6cd0703c-1e23-43c6-96c2-af043e6ad4bf",
	"amount": 1
}
```
