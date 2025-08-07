# 🎰 Transaction management system

Welcome to simple transaction management system for a casino. The system will track user transactions related to their bets and wins. The transactions will be processed asynchronously via a message system, and the data will be stored in a relational database.

### ⚙️ Prerequisites

- [golang](https://go.dev/)
- [MySQL](https://documentation.ubuntu.com/server/how-to/databases/install-mysql/)
- [RabbitMQ](https://www.rabbitmq.com/docs/download)
- set up `.env` with mysql connection string (see `.env.example`) 

RabbitMQ Publisher: Continuously publishes messages with a period of 1 millisecond (it's enough to produce many messages)

RabbitMQ Consumer: Receives, processes, and stores messages in a MySQL database

REST API: Listens on `localhost:8080/transactions` for HTTP requests

Graceful Shutdown: Handles `CTRL+C` signal to cleanly terminate the application

### 🚀 Getting Started

Start the program using just one command!

`make start`

To see all available commands, run:

`make help`

### 💡 Improvements

- Dockerize the application and use containerized MySQL and RabbitMQ
- Write CI/CD workflow to test the application when push occurs
- If number of `make` jobs increases, separate `Makefile` and shell script which will run the commands

### Test Application

To test the application run: 

`make test`

To test sub-packages (e.g consumer) run:

`make test-{pkg}`. (e.g, `make test-consumer`)

To see the coverage percentage for the whole application run:

`make test-cover`

### Coverage

The application coverage is available by:

`make coverage`

To cover sub-package, use:

`make cvr-{pkg}` (e.g, `make cvr-consumer`)

### Test API

API provides transactions filtered by user id and transaction type. Additionally limiter provides last N transaction. 

Get all transactions: 

`curl http://localhost:8080/transactions?`

Get all user's transactions:

`curl http://localhost:8080/transactions?user_id={USER_ID}`

Get all `bet/win` transactions:

`curl http://localhost:8080/transactions?transaction_type={TRANSACTION_TYPE}`

Get last N `bet/win` user's transactions:

`curl http://localhost:8080/transactions?user_id={USER_ID}&transaction_type={TRANSACTION_TYPE}&limit={LIMIT}`
