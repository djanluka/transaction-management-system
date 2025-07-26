### Test Application

To test application run: 

`make test`

To test sub-packages (e.g consumer) run:

`make test-{pkg}`. (e.g, `make test-consumer`)

To see the coverage percentage run:

`make test-cover`

### Coverage

The application coverage is available by:

`make coverage`

To cover sub-package, use:

`make cvr-{pkg}` (e.g, `make cvr-consumer`)

### Test API

API provides transactions filtered by user id and transaction type. Additionally limiter provides last N transaction. 

Get all transactions: `curl http://localhost:8080/transactions?`

Get all user's transactions: `curl http://localhost:8080/transactions?user_id={USER_ID}`

Get all `bet/win` transactions: `curl http://localhost:8080/transactions?transaction_type={TRANSACTION_TYPE}`

Get last N `bet/win` user's transactions: `curl http://localhost:8080/transactions?user_id={USER_ID}&transaction_type={USER_ID}&limit={LIMIT}`
