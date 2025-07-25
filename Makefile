.PHONY: coverage start test clean help

# Default target
.DEFAULT_GOAL := help

# Generate coverage report
coverage:
	@echo "Generating coverage report for $(or $(PKG),all packages)..."
	go test -coverprofile=coverage.out $(or $(PKG),./...)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html

# Run the application
start:
	@echo "Starting application..."
	go run main.go

# Run tests
test:
	@echo "Running tests for all packages"
	ENV_PATH=../.env go test -v -cover ./...
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


# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html

# Show help
help:
	@echo "Available targets:"
	@echo "  make test [package]  - Run tests (all or specific package)"
	@echo "  make coverage [package] - Generate coverage report"
	@echo "  make start        - Run the application"
	@echo "  make clean        - Remove generated files"
	@echo "  make help         - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make test          # Test all packages"
	@echo "  make test consumer # Test consumer package"
	@echo "  make coverage repository # Coverage for repository"