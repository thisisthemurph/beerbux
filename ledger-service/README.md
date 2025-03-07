# ledger-service

The ledger-service listens for `transaction.created` events, validates that ledger, 
and publishes events detailing the ledger with the newly created ledger entries.

## Manual testing

```shell
nats pub transaction.created '{"version": "1.0.0", "data": {"transaction_id": "47218b87-3cbb-49c3-a5fc-62b992f35174", "creator_id": "f3232bbf-c579-4995-9012-a75c7cdec425", "session_id": "94eac316-2900-4e1d-bb56-f67f66309b3c", "member_amounts": [{"user_id": "cae11f2d-f0f8-4b38-8cdb-21c979b5ca48", "amount": 1}, {"user_id": "986dd056-0d15-4148-a935-96b8588aea4c", "amount": 1}]}}'
```
