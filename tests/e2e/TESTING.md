# E2E Testing Documentation

This document describes the end-to-end testing approach for the Price Tracker Backend Service (part of the complete crypto trading platform).

## Overview

The E2E tests are designed to test the entire backend application from outside using HTTP calls. These tests verify the complete flow from server startup to API responses, ensuring the backend service works correctly as a whole. Note that these tests focus on the backend API and do not test the React frontend UI.

## Test Structure

### Test Files
- `e2e_test.go` - Main E2E test file containing all test cases
- `test_e2e.sh` - Shell script for running tests with proper setup/cleanup
- `TESTING.md` - This documentation file

### Test Categories

1. **Health Check Tests**
   - Verify the health endpoint responds correctly
   - Test server startup and readiness

2. **API Endpoint Tests**
   - Test candles endpoint with various parameters
   - Verify error handling for invalid requests
   - Test different symbols and periods

3. **Response Structure Tests**
   - Validate JSON response format
   - Verify OHLC data integrity
   - Check time field consistency

4. **Concurrency Tests**
   - Test multiple simultaneous requests
   - Verify server stability under load

5. **Graceful Shutdown Tests**
   - Test server shutdown behavior
   - Verify cleanup processes

6. **End-to-End Flow Tests**
   - Test complete application workflow
   - Verify multiple symbol/period combinations

## Running Tests

### Prerequisites
- Go 1.21 or higher
- SQLite3
- Bash shell (for test script)

### Quick Start
```bash
# Run all E2E tests with the test script
make e2e-test

# Or run directly
./test_e2e.sh
```

### Manual Test Execution
```bash
# Run tests with verbose output
make e2e-test-verbose

# Run specific test
go test -v -run TestHealthEndpoint ./e2e_test.go

# Run tests with custom timeout
go test -v -timeout=10m ./e2e_test.go
```

### Test Commands
```bash
# Build and run tests
make test-all

# Run only E2E tests
make e2e-test

# Run with verbose output
make e2e-test-verbose
```

## Test Configuration

### Test Environment
- **Base URL**: `http://localhost:8080`
- **Server Port**: 8080
- **Timeout**: 30 seconds per test
- **Database**: Temporary SQLite file (`price_tracker.db`)

### Test Data
- **Symbols**: BTCUSDT, ETHUSDT, ADAUSDT, DOTUSDT, LINKUSDT
- **Periods**: 1, 2, 5, 10 minutes
- **Test Duration**: ~2-3 minutes for full suite

## Test Cases Details

### 1. Health Endpoint Test
```go
func TestHealthEndpoint(t *testing.T)
```
- Starts the server
- Calls `/health` endpoint
- Verifies response status and content
- Checks for "healthy" status

### 2. Candles Endpoint Without Data
```go
func TestCandlesEndpointWithoutData(t *testing.T)
```
- Tests candles endpoint with no historical data
- Verifies empty response structure
- Validates symbol and period parameters

### 3. Invalid Parameters Test
```go
func TestCandlesEndpointInvalidParameters(t *testing.T)
```
- Tests missing required parameters
- Tests invalid parameter values
- Verifies proper error responses (400 Bad Request)

### 4. Different Periods Test
```go
func TestCandlesEndpointWithDifferentPeriods(t *testing.T)
```
- Tests all configured periods (1, 2, 5, 10 minutes)
- Verifies correct period handling
- Validates response structure for each period

### 5. Different Symbols Test
```go
func TestCandlesEndpointWithDifferentSymbols(t *testing.T)
```
- Tests all configured symbols
- Verifies symbol parameter handling
- Validates response for each symbol

### 6. Response Structure Test
```go
func TestCandlesResponseStructure(t *testing.T)
```
- Validates JSON response format
- Checks Content-Type header
- Verifies OHLC data integrity
- Tests time field consistency

### 7. Concurrent Requests Test
```go
func TestConcurrentRequests(t *testing.T)
```
- Sends 10 concurrent requests
- Verifies server stability
- Tests response consistency

### 8. Graceful Shutdown Test
```go
func TestServerGracefulShutdown(t *testing.T)
```
- Tests server shutdown behavior
- Verifies cleanup processes
- Checks for proper exit

### 9. End-to-End Flow Test
```go
func TestEndToEndFlow(t *testing.T)
```
- Tests complete application workflow
- Verifies multiple symbol/period combinations
- Tests health check and candles endpoints together

## Test Utilities

### Server Management
```go
func startServer(t *testing.T) (*exec.Cmd, func())
```
- Builds the application
- Starts the server process
- Returns cleanup function
- Handles process lifecycle

### Server Readiness Check
```go
func waitForServerReady(t *testing.T, baseURL string, timeout time.Duration)
```
- Polls health endpoint
- Waits for server to be ready
- Handles timeout scenarios

### Test Configuration
```go
func setupTestConfig() *TestConfig
```
- Sets up test configuration
- Defines base URL and timeouts
- Configures test parameters

## Response Validation

### Health Response
```json
{
  "status": "healthy"
}
```

### Candles Response
```json
{
  "symbol": "BTCUSDT",
  "period": 1,
  "candles": [
    {
      "symbol": "BTCUSDT",
      "period": 1,
      "open": 50000.0,
      "high": 50100.0,
      "low": 49900.0,
      "close": 50050.0,
      "volume": 100.5,
      "open_time": "2024-01-01T10:00:00Z",
      "close_time": "2024-01-01T10:01:00Z"
    }
  ]
}
```

## Error Handling

### Expected Error Responses
- **400 Bad Request**: Missing or invalid parameters
- **500 Internal Server Error**: Server-side errors
- **Connection Refused**: Server not running

### Test Error Scenarios
- Missing symbol parameter
- Missing period parameter
- Invalid period value
- Empty symbol parameter
- Server not responding

## Performance Considerations

### Test Execution Time
- **Individual Test**: 5-10 seconds
- **Full Suite**: 2-3 minutes
- **Concurrent Tests**: 10 simultaneous requests

### Resource Usage
- **Memory**: ~50MB per test run
- **CPU**: Minimal during idle
- **Disk**: Temporary SQLite database

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check for running processes
   lsof -i :8080
   
   # Kill existing processes
   pkill -f "price-tracker"
   ```

2. **Test Timeout**
   ```bash
   # Increase timeout
   go test -v -timeout=10m ./e2e_test.go
   ```

3. **Build Failures**
   ```bash
   # Clean and rebuild
   make clean
   go mod tidy
   go build
   ```

4. **Database Lock Issues**
   ```bash
   # Remove database file
   rm -f price_tracker.db
   ```

### Debug Mode
```bash
# Run with debug output
go test -v -timeout=5m ./e2e_test.go -debug

# Run specific test with verbose output
go test -v -run TestHealthEndpoint ./e2e_test.go
```

## Continuous Integration

### GitHub Actions Example
```yaml
name: E2E Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      - run: make e2e-test
```

### Local CI
```bash
# Run full test suite
make test-all

# Run with coverage (when unit tests are added)
go test -cover ./...
```

## Best Practices

1. **Isolation**: Each test runs in isolation with fresh server instance
2. **Cleanup**: Proper cleanup of resources after each test
3. **Timeouts**: Appropriate timeouts for all operations
4. **Validation**: Comprehensive response validation
5. **Error Handling**: Test both success and error scenarios
6. **Concurrency**: Test server stability under load

## Future Enhancements

1. **Mock WebSocket Data**: Add tests with simulated trade data
2. **Performance Tests**: Add load testing scenarios
3. **Integration Tests**: Test with real Binance WebSocket
4. **Coverage Reports**: Add test coverage metrics
5. **Parallel Execution**: Run tests in parallel where possible 