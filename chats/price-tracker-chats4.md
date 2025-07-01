# Build go backend for Binance WebSocket
_Exported on 01/07/2025 at 14:20:33 GMT+4 from Cursor (1.1.6)_

---

**User**

Write me a React JS UI and “go” back end service.

The back end service will have those below :


After retrieving which symbols to listen to and the URL to connect to from the YAML configuration file, it will start listening to the WebSocket stream.
https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams

We need to subscribe to <symbol>@trade stream only. For now you can ignore other event types.

The backend will generate OHLC (Open, High, Low, Close) from the incoming prices. While doing this, it will save the OHLC data it produces in different periods like 1 minute, 2 minutes, 5 minutes, and 10 minutes into SQLite. It will respond with list of candle sticks created from the trades in the last hour only, from to requests the interface for a specific symbol and period. Be careful about the time zone differences coming from the stream. You need to pick the last hour done with the latest trade in database.For now only write the back end service. 
Do not write any tests. I will give extra instructions to mock the stream with wire mock or  a test container later.

---

**Cursor**

I'm currently using Claude Sonnet 4, an AI model designed to help with coding tasks.

I'll create a React JS UI and Go backend service for a price tracker that connects to Binance WebSocket streams and generates OHLC data. Let me start by exploring the workspace structure and then building the backend service.


The workspace is empty, so I'll create the complete project structure. Let me start with the Go backend service following SOLID, KISS, DRY, and YAGNI principles.


```go.mod
module price-tracker

go 1.21

require (
	github.com/gorilla/websocket v1.5.1
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/gin-gonic/gin v1.9.1
	github.com/gin-contrib/cors v1.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
```

```yaml
websocket:
  url: "wss://stream.binance.com:9443/ws"
  symbols:
    - "btcusdt"
    - "ethusdt"
    - "adausdt"
    - "dogeusdt"

server:
  port: 8080
  host: "localhost"

database:
  path: "./price_tracker.db"
```

---

**User**

On Live Trades pages the table is frequently updated and that causes blinks and shakes..

---

**Cursor**

I see you've made some changes to the configuration and are experiencing UI issues with the Live Trades page. Let me first complete the backend service and then address the UI blinking/shaking issue.

Let me continue building the Go backend service with the updated configuration:


```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"price-tracker/internal/config"
	"price-tracker/internal/database"
	"price-tracker/internal/handler"
	"price-tracker/internal/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize WebSocket client
	wsClient := websocket.NewClient(cfg.WebSocket.URL, cfg.WebSocket.Symbols, db)

	// Start WebSocket connection
	go func() {
		if err := wsClient.Start(); err != nil {
			log.Printf("WebSocket error: %v", err)
		}
	}()

	// Initialize HTTP handler
	h := handler.New(db)

	// Start HTTP server
	go func() {
		if err := h.Start(cfg.Server.Host, cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s:%d", cfg.Server.Host, cfg.Server.Port)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	wsClient.Stop()
}
```

I need to create the internal packages first. Let me create the directory structure and the required packages:


```go
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WebSocket WebSocketConfig `yaml:"websocket"`
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	OHLCPeriods []int         `yaml:"ohlc_periods"`
}

type WebSocketConfig struct {
	URL     string   `yaml:"url"`
	Symbols []string `yaml:"symbols"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
```

```go
package models

import (
	"strconv"
	"time"
)

type Trade struct {
	Symbol          string  `json:"s"`
	ID              int64   `json:"t"`
	Price           string  `json:"p"`
	Quantity        string  `json:"q"`
	BuyerOrderID    int64   `json:"b"`
	SellerOrderID   int64   `json:"a"`
	TradeTime       int64   `json:"T"`
	IsBuyerMaker    bool    `json:"m"`
	Ignore          bool    `json:"M"`
	ParsedPrice     float64
	ParsedQuantity  float64
	ParsedTime      time.Time
}

func (t *Trade) Parse() error {
	price, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		return err
	}
	t.ParsedPrice = price

	quantity, err := strconv.ParseFloat(t.Quantity, 64)
	if err != nil {
		return err
	}
	t.ParsedQuantity = quantity

	// Convert milliseconds to time.Time
	t.ParsedTime = time.Unix(0, t.TradeTime*int64(time.Millisecond))
	return nil
}
```

```go
package models

import (
	"time"
)

type OHLC struct {
	ID        int64     `json:"id"`
	Symbol    string    `json:"symbol"`
	Period    int       `json:"period"` // in minutes
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

type OHLCRequest struct {
	Symbol string `json:"symbol" form:"symbol" binding:"required"`
	Period int    `json:"period" form:"period" binding:"required"`
}

type OHLCResponse struct {
	Symbol    string `json:"symbol"`
	Period    int    `json:"period"`
	Candles   []OHLC `json:"candles"`
	Count     int    `json:"count"`
	FromTime  string `json:"from_time"`
	ToTime    string `json:"to_time"`
}
```

```go
package database

import (
	"database/sql"
	"fmt"
	"time"

	"price-tracker/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func New(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ohlc (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		period INTEGER NOT NULL,
		open REAL NOT NULL,
		high REAL NOT NULL,
		low REAL NOT NULL,
		close REAL NOT NULL,
		volume REAL NOT NULL,
		timestamp DATETIME NOT NULL,
		UNIQUE(symbol, period, timestamp)
	);
	CREATE INDEX IF NOT EXISTS idx_oh
```

