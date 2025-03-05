# transaction-service

The transaction-service receives gRPC requests to create a transaction, validates that transaction, and publishes an
event detailing the transaction with the new transaction ID.

**Create a transaction**

```shell
grpcurl -plaintext -d '{"creator_id": "0776feab-7b57-4be0-b5b8-f57772d572d3", "session_id": "9261d60b-0de6-45ff-95f1-6cb3e056c05f", "member_amounts": [{"user_id": "90627371-879e-4fcd-8100-cad42f46bfa9", "amount": 1}]}' localhost:50053 Transaction.CreateTransaction
```
