# Tests Directory

This directory contains all test files for the Price Tracker - Complete Crypto Trading Platform.

## Directory Structure

```
tests/
├── README.md           # This file
├── e2e/                # End-to-End tests
│   ├── e2e_test.go     # Main E2E test file
│   ├── test_e2e.sh     # E2E test runner script
│   └── TESTING.md      # E2E testing documentation
└── unit/               # Unit tests (future)
    └── (unit test files will go here)
```

## Test Types

### E2E Tests (`tests/e2e/`)
- **Purpose**: Test the entire backend application from outside using HTTP calls
- **Scope**: Full backend application integration testing
- **Files**: 
  - `e2e_test.go` - Main test cases
  - `test_e2e.sh` - Test runner script
  - `TESTING.md` - Detailed documentation

### Unit Tests (`tests/unit/`) - Future
- **Purpose**: Test individual backend components in isolation
- **Scope**: Package-level testing
- **Files**: Will be added when unit tests are implemented

## Running Tests

### E2E Tests
```bash
# From project root
make e2e-test

# Or run directly
./tests/e2e/test_e2e.sh

# Verbose output
make e2e-test-verbose
```

### All Tests
```bash
# Run all tests (when unit tests are added)
make test-all
```

## Test Organization Principles

1. **Separation of Concerns**: E2E and unit tests are in separate directories
2. **Self-Contained**: Each test type has its own runner and documentation
3. **Easy Discovery**: Clear directory structure makes it easy to find tests
4. **Maintainable**: Tests are organized by type and scope

## Adding New Tests

### For E2E Tests
1. Add test functions to `tests/e2e/e2e_test.go`
2. Update `tests/e2e/TESTING.md` if needed
3. Test with `make e2e-test`

### For Unit Tests (Future)
1. Create test files in `tests/unit/`
2. Follow Go testing conventions
3. Run with `go test ./tests/unit/`

## Test Configuration

- **E2E Tests**: Use project root for building and running
- **Unit Tests**: Will use package-level testing
- **Test Data**: Stored in temporary files that are cleaned up automatically
- **Dependencies**: Managed through Go modules

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Always clean up resources after tests
3. **Documentation**: Keep test documentation up to date
4. **Naming**: Use descriptive test names and file names
5. **Organization**: Group related tests together 