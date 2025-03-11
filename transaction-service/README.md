# transaction-service

The transaction-service receives gRPC requests to create a transaction, validates that transaction, and publishes an
event detailing the transaction with the new transaction ID.

**Create a transaction**

```shell
grpcurl -plaintext -d '{"creator_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde", "session_id": "5c0327eb-b934-46be-a882-56195fab04d9", "member_amounts": [{"user_id": "6cd0703c-1e23-43c6-96c2-af043e6ad4bf", "amount": 1}]}' localhost:50053 Transaction.CreateTransaction
```
