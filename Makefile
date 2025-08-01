.PHONY: coverage start test clean help

# Default target
.DEFAULT_GOAL := help

# Generate coverage report
cvr:
	@echo "Generating coverage report for packages"
	ENV_PATH=../.env go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html
cvr-consumer:
	@echo "Generating coverage report for consumer"
	ENV_PATH=../.env go test -coverprofile=coverage.out ./consumer
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html
cvr-database:
	@echo "Generating coverage report for database"
	ENV_PATH=../.env go test -coverprofile=coverage.out ./database
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html
cvr-publisher:
	@echo "Generating coverage report for publisher"
	ENV_PATH=../.env go test -coverprofile=coverage.out ./publisher
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html
cvr-rabbitmq:
	@echo "Generating coverage report for rabbitmq"
	ENV_PATH=../.env go test -coverprofile=coverage.out ./rabbitmq
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html

# This cover test fails due to 'signal: terminated' that is used for graceful shutdown test
# 'make test-transaction' passes well
# cvr-transaction:
# 	@echo "Generating coverage report for transaction"
# 	ENV_PATH=../.env go test -coverprofile=coverage.out ./transaction
# 	go tool cover -html=coverage.out -o coverage.html
# 	@echo "Opening coverage report..."
# 	open coverage.html

# Run the application
start:
	@echo "Pre-init database"
	mysql < database/migrations/init.sql
	@echo "Starting application..."
	go run main.go

# Run tests
test:
	@echo "Running tests for all packages"
	ENV_PATH=../.env go test -v -cover ./...
test-cover:
	@echo "Running tests for all packages"
	ENV_PATH=../.env go test -cover ./...
test-consumer:
	@echo "Running tests for consumer"
	ENV_PATH=../.env go test -v -cover ./consumer
test-database:
	@echo "Running tests for database"
	ENV_PATH=../.env go test -v -cover ./database
test-publisher:
	@echo "Running tests for publisher"
	ENV_PATH=../.env go test -v -cover ./publisher
test-rabbitmq:
	@echo "Running tests for rabbitmq"
	ENV_PATH=../.env go test -v -cover ./rabbitmq
test-transaction:
	@echo "Running tests for transaction"
	ENV_PATH=../.env go test -v -cover ./transaction


# Reset the database table
reset:
	@echo "Reset mysql casino table"
	mysql < database/migrations/reset.sql


# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html

# Show help
help:
	@echo "Available targets:"
	@echo "  make test 				- Run tests"
	@echo "  make test-cover 			- Run tests and print coverage percentage"
	@echo "  make test-{pkg}			- Run tests for specific package (consumer/database/publisher/etc)"
	@echo "  make cvr			 	- Generate coverage report"
	@echo "  make cvr-{pkg}			- Generate coverage report for specific package(consumer/database/publisher/etc)"
	@echo "  make start        			- Run the application"
	@echo "  make reset        			- Reset the database"
	@echo "  make clean        			- Remove generated files"
	@echo "  make help         			- Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make test 		# Test all packages"
	@echo "  make test-consumer	# Test consumer package"
	@echo "  make start		# Start the application"