package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type TestConfig struct {
	BaseURL    string
	ServerPort int
	Timeout    time.Duration
}

type CandleResponse struct {
	Symbol  string        `json:"symbol"`
	Period  int           `json:"period"`
	Candles []CandleStick `json:"candles"`
}

type CandleStick struct {
	Symbol    string    `json:"symbol"`
	Period    int       `json:"period"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type MockTrade struct {
	Symbol    string `json:"s"`
	TradeID   int64  `json:"t"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeTime int64  `json:"T"`
}

func setupTestConfig() *TestConfig {
	return &TestConfig{
		BaseURL:    "http://localhost:8080",
		ServerPort: 8080,
		Timeout:    30 * time.Second,
	}
}

func startServer(t *testing.T) (*exec.Cmd, func()) {
	// Get the project root directory (two levels up from tests/e2e)
	projectRoot, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Build the application from project root
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(projectRoot, "test-price-tracker"), filepath.Join(projectRoot, "main.go"))
	buildCmd.Dir = projectRoot
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}

	// Start the server from project root
	cmd := exec.Command(filepath.Join(projectRoot, "test-price-tracker"))
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to start
	time.Sleep(3 * time.Second)

	cleanup := func() {
		cmd.Process.Kill()
		cmd.Wait()
		os.Remove(filepath.Join(projectRoot, "test-price-tracker"))
		os.Remove(filepath.Join(projectRoot, "price_tracker.db"))
	}

	return cmd, cleanup
}

func waitForServerReady(t *testing.T, baseURL string, timeout time.Duration) {
	client := &http.Client{Timeout: timeout}
	start := time.Now()

	for time.Since(start) < timeout {
		resp, err := client.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("Server not ready after %v", timeout)
}

func TestHealthEndpoint(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	// Test health endpoint
	resp, err := http.Get(config.BaseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", health.Status)
	}
}

func TestCandlesEndpointWithoutData(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	// Test candles endpoint without any data
	resp, err := http.Get(config.BaseURL + "/api/candles?symbol=BTCUSDT&period=1")
	if err != nil {
		t.Fatalf("Failed to call candles endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var candles CandleResponse
	if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
		t.Fatalf("Failed to decode candles response: %v", err)
	}

	if candles.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", candles.Symbol)
	}

	if candles.Period != 1 {
		t.Errorf("Expected period 1, got %d", candles.Period)
	}

	// Should return empty candles array initially
	if len(candles.Candles) != 0 {
		t.Errorf("Expected empty candles array, got %d candles", len(candles.Candles))
	}
}

func TestCandlesEndpointInvalidParameters(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	testCases := []struct {
		name     string
		url      string
		expected int
	}{
		{"Missing symbol", "/api/candles?period=1", http.StatusBadRequest},
		{"Missing period", "/api/candles?symbol=BTCUSDT", http.StatusBadRequest},
		{"Invalid period", "/api/candles?symbol=BTCUSDT&period=invalid", http.StatusBadRequest},
		{"Empty symbol", "/api/candles?symbol=&period=1", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(config.BaseURL + tc.url)
			if err != nil {
				t.Fatalf("Failed to call candles endpoint: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expected {
				t.Errorf("Expected status %d, got %d", tc.expected, resp.StatusCode)
			}
		})
	}
}

func TestCandlesEndpointWithDifferentPeriods(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	periods := []int{1, 2, 5, 10}

	for _, period := range periods {
		t.Run(fmt.Sprintf("Period_%d", period), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/api/candles?symbol=ETHUSDT&period=%d", config.BaseURL, period))
			if err != nil {
				t.Fatalf("Failed to call candles endpoint: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}

			var candles CandleResponse
			if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
				t.Fatalf("Failed to decode candles response: %v", err)
			}

			if candles.Symbol != "ETHUSDT" {
				t.Errorf("Expected symbol 'ETHUSDT', got '%s'", candles.Symbol)
			}

			if candles.Period != period {
				t.Errorf("Expected period %d, got %d", period, candles.Period)
			}
		})
	}
}

func TestCandlesEndpointWithDifferentSymbols(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "DOTUSDT", "LINKUSDT"}

	for _, symbol := range symbols {
		t.Run(fmt.Sprintf("Symbol_%s", symbol), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/api/candles?symbol=%s&period=1", config.BaseURL, symbol))
			if err != nil {
				t.Fatalf("Failed to call candles endpoint: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}

			var candles CandleResponse
			if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
				t.Fatalf("Failed to decode candles response: %v", err)
			}

			if candles.Symbol != symbol {
				t.Errorf("Expected symbol '%s', got '%s'", symbol, candles.Symbol)
			}

			if candles.Period != 1 {
				t.Errorf("Expected period 1, got %d", candles.Period)
			}
		})
	}
}

func TestCandlesResponseStructure(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	resp, err := http.Get(config.BaseURL + "/api/candles?symbol=BTCUSDT&period=5")
	if err != nil {
		t.Fatalf("Failed to call candles endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain 'application/json', got '%s'", contentType)
	}

	var candles CandleResponse
	if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
		t.Fatalf("Failed to decode candles response: %v", err)
	}

	// Verify response structure
	if candles.Symbol == "" {
		t.Error("Expected symbol to be present in response")
	}

	if candles.Period <= 0 {
		t.Error("Expected period to be positive")
	}

	// If there are candles, verify their structure
	for i, candle := range candles.Candles {
		if candle.Symbol == "" {
			t.Errorf("Candle %d: Expected symbol to be present", i)
		}

		if candle.Period <= 0 {
			t.Errorf("Candle %d: Expected period to be positive", i)
		}

		if candle.Open < 0 {
			t.Errorf("Candle %d: Expected open price to be non-negative", i)
		}

		if candle.High < 0 {
			t.Errorf("Candle %d: Expected high price to be non-negative", i)
		}

		if candle.Low < 0 {
			t.Errorf("Candle %d: Expected low price to be non-negative", i)
		}

		if candle.Close < 0 {
			t.Errorf("Candle %d: Expected close price to be non-negative", i)
		}

		if candle.Volume < 0 {
			t.Errorf("Candle %d: Expected volume to be non-negative", i)
		}

		// Verify OHLC logic
		if candle.High < candle.Low {
			t.Errorf("Candle %d: High price cannot be less than low price", i)
		}

		if candle.High < candle.Open || candle.High < candle.Close {
			t.Errorf("Candle %d: High price should be >= open and close", i)
		}

		if candle.Low > candle.Open || candle.Low > candle.Close {
			t.Errorf("Candle %d: Low price should be <= open and close", i)
		}

		// Verify time fields
		if candle.OpenTime.IsZero() {
			t.Errorf("Candle %d: Expected open_time to be set", i)
		}

		if candle.CloseTime.IsZero() {
			t.Errorf("Candle %d: Expected close_time to be set", i)
		}

		if candle.CloseTime.Before(candle.OpenTime) {
			t.Errorf("Candle %d: Close time should be after open time", i)
		}
	}
}

func TestConcurrentRequests(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	// Test concurrent requests to the same endpoint
	numRequests := 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := http.Get(config.BaseURL + "/api/candles?symbol=BTCUSDT&period=1")
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				return
			}

			var candles CandleResponse
			if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
				results <- err
				return
			}

			results <- nil
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent request failed: %v", err)
		}
	}
}

func TestServerGracefulShutdown(t *testing.T) {
	config := setupTestConfig()
	cmd, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	// Verify server is running
	resp, err := http.Get(config.BaseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Send SIGTERM to gracefully shutdown
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("Failed to send interrupt signal: %v", err)
	}

	// Wait for process to exit
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Process exited with error: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Error("Process did not exit within 10 seconds")
	}
}

func TestEndToEndFlow(t *testing.T) {
	config := setupTestConfig()
	_, cleanup := startServer(t)
	defer cleanup()

	waitForServerReady(t, config.BaseURL, config.Timeout)

	// Test the complete flow
	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(config.BaseURL + "/health")
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Health check failed with status %d", resp.StatusCode)
		}
	})

	t.Run("InitialCandles", func(t *testing.T) {
		resp, err := http.Get(config.BaseURL + "/api/candles?symbol=BTCUSDT&period=1")
		if err != nil {
			t.Fatalf("Failed to get initial candles: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to get initial candles with status %d", resp.StatusCode)
		}

		var candles CandleResponse
		if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
			t.Fatalf("Failed to decode candles response: %v", err)
		}

		// Initially should have no candles
		if len(candles.Candles) != 0 {
			t.Errorf("Expected no candles initially, got %d", len(candles.Candles))
		}
	})

	t.Run("MultipleSymbols", func(t *testing.T) {
		symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"}
		periods := []int{1, 2, 5}

		for _, symbol := range symbols {
			for _, period := range periods {
				resp, err := http.Get(fmt.Sprintf("%s/api/candles?symbol=%s&period=%d", config.BaseURL, symbol, period))
				if err != nil {
					t.Errorf("Failed to get candles for %s period %d: %v", symbol, period, err)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					t.Errorf("Failed to get candles for %s period %d with status %d", symbol, period, resp.StatusCode)
					continue
				}

				var candles CandleResponse
				if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
					t.Errorf("Failed to decode candles for %s period %d: %v", symbol, period, err)
					continue
				}

				if candles.Symbol != symbol {
					t.Errorf("Expected symbol %s, got %s", symbol, candles.Symbol)
				}

				if candles.Period != period {
					t.Errorf("Expected period %d, got %d", period, candles.Period)
				}
			}
		}
	})
}
