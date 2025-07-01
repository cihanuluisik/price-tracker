# Unified Chat Requirements: Crypto Price Tracker Backend & UI

## **Initial Backend Requirements**

### **Core Backend Service (Go)**
Write me a "go" back end service with the following features:

1. **Configuration Management**
   - Crypto symbols provided in YAML config file
   - Server configuration (host, port)
   - Database configuration (SQLite path)
   - OHLC periods configuration

2. **WebSocket Integration**
   - Connect to Binance WebSocket streams
   - Subscribe to `<symbol>@trade` stream only
   - Handle real-time trade data
   - Parse trade information (price, quantity, timestamp)

3. **OHLC Generation**
   - Generate OHLC (Open, High, Low, Close) from incoming prices
   - Support multiple periods: 10s, 30s, 1m, 2m, 5m
   - Save OHLC data to SQLite database
   - Handle timezone differences from stream

4. **API Endpoints**
   - Health check endpoint (`/health`)
   - Candles endpoint (`/api/candles`) with parameters:
     - `symbol` (required)
     - `period` (required)
     - `hours` (optional, default 1)
   - Return candlestick data for the last hour
   - Include trade count and latest trade information

5. **Database Design**
   - SQLite database with trades table
   - Store individual trades with parsed data
   - Generate candlesticks dynamically from trades table
   - Support for both seconds and minutes periods

## **Testing Requirements**

### **E2E Testing**
- Write end-to-end tests using HTTP calls (no unit tests)
- Test the complete application flow
- Include health endpoint tests
- Test candles endpoint with various parameters
- Mock WebSocket streams for testing
- Move tests to proper directory structure (`tests/e2e/`)

## **UI Requirements**

### **React JS Frontend**
Build a React JS UI with the following features:

1. **Live Trades Page**
   - Real-time trade data display
   - WebSocket connection to backend
   - Table with trade information
   - Fix blinking/shaking issues with frequent updates
   - Connection status indicators

2. **Price Charts Page**
   - Interactive candlestick charts
   - Symbol selector (BTCUSDT, ETHUSDT, ADAUSDT, DOTUSDT, LINKUSDT)
   - Period selector with new format:
     - 10s (default)
     - 30s
     - 1m
     - 2m
     - 5m
   - Real-time chart updates
   - Chart controls and information display

3. **UI Components**
   - Modern, responsive design
   - Navigation between pages
   - Loading states and error handling
   - Connection status indicators

## **Configuration Updates**

### **YAML Configuration**
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

## **Development & Deployment**

### **Project Structure**
```
price-tracker-1/
├── api/           # HTTP handlers
├── config/        # Configuration management
├── database/      # Database operations
├── internal/      # Internal packages
├── models/        # Data models
├── ohlc/          # OHLC generation
├── tests/         # Test files
│   └── e2e/       # End-to-end tests
├── ui/            # React frontend
├── websocket/     # WebSocket handling
├── config.yaml    # Configuration file
├── go.mod         # Go module
├── main.go        # Main application
├── start.sh       # Startup script
└── .gitignore     # Git ignore patterns
```

### **Startup Script**
- Create `start.sh` script in project root
- Start backend first, then UI
- Handle port conflicts
- Proper cleanup on exit
- Dependency checking (Go, Node.js)
- Build and run both services

### **Code Quality**
- Follow SOLID, KISS, DRY, YAGNI principles
- Create short classes and methods
- No unit tests (only E2E tests)
- Suggest integration tests for classes/methods over 30 lines

## **Key Features to Implement**

1. **Dynamic Candlestick Generation**
   - Generate OHLC data on-the-fly from trades table
   - Support for both seconds and minutes periods
   - Time-based queries instead of limit-based

2. **Real-time Updates**
   - WebSocket connections for live data
   - Efficient UI updates without blinking
   - Connection status management

3. **Error Handling**
   - Graceful WebSocket reconnection
   - API error responses
   - UI error states and loading indicators

4. **Performance Optimization**
   - Efficient database queries
   - Optimized UI rendering
   - Memory management for real-time data

## **Final Deliverables**

1. **Complete Go Backend**
   - WebSocket client for Binance
   - HTTP API server
   - SQLite database with dynamic OHLC generation
   - Configuration management

2. **React Frontend**
   - Live trades page with real-time updates
   - Price charts page with interactive candlesticks
   - Period selection (10s, 30s, 1m, 2m, 5m)
   - Modern UI with responsive design

3. **Testing & Documentation**
   - E2E tests for complete application flow
   - Startup script for easy deployment
   - Proper project structure and organization

4. **Configuration**
   - YAML configuration file
   - Environment-specific settings
   - Proper .gitignore patterns

This unified chat represents all the requirements extracted from the multiple chat sessions, providing a comprehensive guide for building the complete crypto price tracker application. 