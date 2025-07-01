.PHONY: build run clean test e2e-test e2e-test-verbose

# Build the application
build:
	go build -o price-tracker main.go

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -f price-tracker
	rm -f price_tracker.db
	rm -f test-price-tracker

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run unit tests (when added later)
test:
	go test ./...

# Run E2E tests
e2e-test:
	./tests/e2e/test_e2e.sh

# Run E2E tests with verbose output
e2e-test-verbose:
	go test -v -timeout=5m ./tests/e2e/e2e_test.go

# Build and run
dev: build run

# Run all tests
test-all: test e2e-test 