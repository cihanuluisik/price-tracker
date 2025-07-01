# Price Tracker - Complete Crypto Trading Platform

A comprehensive cryptocurrency trading platform with a Go backend service that tracks real-time prices from Binance WebSocket streams, generates OHLC data, and a modern React frontend for live trading visualization.

## 🛠️ Tech Stack

### **Backend (Go)**
- **Language**: Go 1.21
- **Web Framework**: Standard `net/http` package
- **WebSocket**: `github.com/gorilla/websocket` v1.5.1
- **Database**: SQLite with `github.com/mattn/go-sqlite3` v1.14.17
- **Configuration**: YAML with `gopkg.in/yaml.v3` v3.0.1
- **Architecture**: Modular design with separate packages for API, config, database, models, and websocket

### **Frontend (React)**
- **Framework**: React 18.2.0
- **Build Tool**: Create React App (react-scripts 5.0.1)
- **HTTP Client**: Axios v1.6.0
- **Charts**: 
  - Lightweight Charts v4.1.3 (for candlestick charts)
  - Recharts v2.8.0 (for additional charting)
- **Testing**: 
  - Jest DOM v5.17.0
  - React Testing Library v13.4.0
  - User Event v13.5.0
- **WebSocket**: Native WebSocket API
- **Proxy**: Configured to proxy API calls to `http://localhost:8080`

### **Development Tools**
- **Package Manager**: npm (for React), go modules (for Go)
- **Version Control**: Git
- **Database**: SQLite (file-based)
- **Real-time Communication**: WebSocket connections

## 🚀 Features

### Backend Service
- Real-time WebSocket connection to Binance trade streams
- OHLC data generation for multiple time periods (10s, 30s, 1m, 2m, 5m)
- Dynamic candlestick generation from trades table
- SQLite database storage for historical data
- REST API for retrieving candlestick data
- WebSocket server for broadcasting live trades to frontend
- Configurable symbols and time periods via YAML configuration
- Graceful shutdown handling

### Frontend UI
- Modern React.js application with ES8+ JavaScript
- Live Trades page with real-time WebSocket data
- Interactive Price Charts with candlestick visualization
- Period selection (10s, 30s, 1m, 2m, 5m) with 10s as default
- Pages open in new browser windows
- Responsive design for all device sizes
- Auto-reconnection for WebSocket connections
- Beautiful, modern UI with smooth animations

### Testing & Development
- Comprehensive end-to-end (E2E) tests
- Automated test runner scripts
- Easy startup script for both backend and frontend
- Makefile commands for common operations

## 📁 Project Structure

```
price-tracker-1/
├── main.go                 # Backend application entry point
├── go.mod                  # Go module dependencies
├── config.yaml             # Backend configuration file
├── start.sh                # Startup script for backend + frontend
├── Makefile                # Build and test commands
├── README.md               # This file
├── config/
│   └── config.go           # Configuration loading
├── models/
│   └── trade.go            # Data models
├── database/
│   └── sqlite.go           # Database operations
├── ohlc/
│   └── generator.go        # OHLC data generation
├── websocket/
│   ├── binance.go          # Binance WebSocket client
│   └── server.go           # WebSocket server for frontend
├── api/
│   └── handlers.go         # HTTP API handlers
├── ui/                     # React frontend application
│   ├── package.json        # Frontend dependencies
│   ├── src/
│   │   ├── components/
│   │   │   ├── LiveTradesPage.js      # Live trades component
│   │   │   ├── LiveTradesPage.css     # Live trades styles
│   │   │   ├── PriceChartsPage.js     # Price charts component
│   │   │   └── PriceChartsPage.css    # Price charts styles
│   │   ├── App.js          # Main application component
│   │   └── index.js        # Application entry point
│   └── README.md           # Frontend documentation
└── tests/                  # Test suite
    ├── e2e/                # End-to-end tests
    │   ├── e2e_test.go     # E2E test cases
    │   ├── test_e2e.sh     # E2E test runner
    │   └── TESTING.md      # E2E testing documentation
    └── README.md           # Testing documentation
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- Node.js 16+ and npm
- SQLite3

## 🚀 Quick Start

### Option 1: Use the Startup Script (Recommended)
```bash
# Clone the repository
git clone <repository-url>
cd price-tracker-1

# Make the startup script executable and run it
chmod +x start.sh
./start.sh
```

This will:
1. Install Go dependencies
2. Build and start the backend service
3. Install frontend dependencies
4. Start the React development server
5. Open your browser to the application

### Option 2: Manual Setup

#### Backend Setup
```bash
# Install Go dependencies
go mod tidy

# Run the backend service
go run main.go
```

The backend will start on `http://localhost:8080`

#### Frontend Setup
```bash
# Navigate to UI directory
cd ui

# Install dependencies
npm install

# Start the development server
npm start
```

The frontend will start on `http://localhost:3000`

## ⚙️ Configuration

Edit `config.yaml` to customize the backend service:

```yaml
symbols:
  - BTCUSDT
  - ETHUSDT
  - ADAUSDT
  - DOTUSDT
  - LINKUSDT

server:
  port: 8080
  host: "0.0.0.0"

database:
  path: "./price_tracker.db"

ohlc_periods:
  - "10s"
  - "30s"
  - "1m"
  - "2m"
  - "5m"
```

## 🌐 API Endpoints

### Get Candlesticks
```
GET /api/candles?symbol=BTCUSDT&period=1
```

Parameters:
- `symbol` (required): Trading pair symbol (e.g., BTCUSDT)
- `period` (required): Time period (10s, 30s, 1m, 2m, 5m)
- `hours` (optional): Number of hours to fetch (default: 1)

Response:
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

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "healthy"
}
```

### WebSocket Endpoint
```
WebSocket /ws/trades
```

Broadcasts real-time trade data to connected frontend clients.

## 🎨 Frontend Features

### Live Trades Page
- Real-time trade table with symbol, price, quantity, time, and trade ID
- Color-coded price changes (green for increases, red for decreases)
- WebSocket connection status indicator
- Auto-reconnection on disconnection
- Responsive design for all screen sizes

### Price Charts Page
- Interactive candlestick charts using TradingView's lightweight-charts
- Symbol selection (BTC, ETH, ADA, DOT, LINK)
- Time period selection (10s, 30s, 1m, 2m, 5m) with 10s as default
- Real-time chart updates
- Professional trading interface
- Pages open in new browser windows

## 🧪 Testing

### Run E2E Tests
```bash
# Run all E2E tests
make e2e-test

# Run with verbose output
make e2e-test-verbose

# Run tests manually
./tests/e2e/test_e2e.sh
```

### Available Make Commands
```bash
# Build the backend
make build

# Run E2E tests
make e2e-test

# Run tests with verbose output
make e2e-test-verbose

# Clean build artifacts
make clean
```

## 🗄️ Database Schema

The service uses SQLite with the following schema:

```sql
CREATE TABLE ohlc (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL,
    period INTEGER NOT NULL,
    open REAL NOT NULL,
    high REAL NOT NULL,
    low REAL NOT NULL,
    close REAL NOT NULL,
    volume REAL NOT NULL,
    open_time DATETIME NOT NULL,
    close_time DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol, period, open_time)
);
```

## 🔄 How It Works

1. **Configuration Loading**: The backend loads symbols and settings from `config.yaml`
2. **WebSocket Connection**: Connects to Binance WebSocket streams for specified symbols
3. **Trade Processing**: Receives real-time trade data and processes it for OHLC generation
4. **OHLC Generation**: Creates candlestick data for multiple time periods (10s, 30s, 1m, 2m, 5m)
5. **Data Storage**: Saves OHLC data to SQLite database
6. **WebSocket Broadcasting**: Broadcasts live trades to connected frontend clients
7. **API Serving**: Provides REST endpoints for retrieving historical candlestick data
8. **Frontend Display**: React UI displays live trades and interactive charts

## 🎯 Development Principles

The code follows SOLID principles with:
- **Single Responsibility**: Each package has a specific purpose
- **Open/Closed**: Easy to extend with new features
- **Liskov Substitution**: Interfaces are properly defined
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Dependencies are injected

Additional principles:
- **KISS**: Keep It Simple, Stupid
- **DRY**: Don't Repeat Yourself
- **YAGNI**: You Aren't Gonna Need It

## 🚀 Future Enhancements

- Add authentication and authorization
- Implement rate limiting
- Add metrics and monitoring
- Support for additional exchanges
- Real-time alerts and notifications
- Data aggregation for longer time periods
- Mobile app development
- Advanced charting features
- Trading bot integration 