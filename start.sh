#!/bin/bash

# Price Tracker Startup Script
# Starts the Go backend first, then the React UI

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Function to check if a port is in use
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Function to wait for server to be ready
wait_for_server() {
    local port=$1
    local max_attempts=30
    local attempt=1
    
    print_info "Waiting for server to be ready on port $port..."
    
    while [ $attempt -le $max_attempts ]; do
        if check_port $port; then
            # Try to connect to health endpoint
            if curl -s http://localhost:$port/health >/dev/null 2>&1; then
                print_status "Server is ready on port $port!"
                return 0
            fi
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    print_error "Server failed to start on port $port after $max_attempts attempts"
    return 1
}

# Function to cleanup on exit
cleanup() {
    print_info "Shutting down services..."
    
    # Kill backend process
    if [ ! -z "$BACKEND_PID" ]; then
        print_info "Stopping backend (PID: $BACKEND_PID)..."
        kill $BACKEND_PID 2>/dev/null || true
    fi
    
    # Kill UI process
    if [ ! -z "$UI_PID" ]; then
        print_info "Stopping UI (PID: $UI_PID)..."
        kill $UI_PID 2>/dev/null || true
    fi
    
    # Kill any remaining processes
    pkill -f "price-tracker" 2>/dev/null || true
    pkill -f "react-scripts" 2>/dev/null || true
    
    print_status "Cleanup completed"
}

# Set up trap to ensure cleanup on exit
trap cleanup EXIT

echo "🚀 Starting Price Tracker Application"
echo "======================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go first."
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    print_error "main.go not found. Please run this script from the project root."
    exit 1
fi

if [ ! -d "ui" ]; then
    print_error "ui directory not found. Please ensure the UI project exists."
    exit 1
fi

print_status "Go version: $(go version)"
print_status "Node.js version: $(node --version)"

# Check if ports are already in use
if check_port 8080; then
    print_warning "Port 8080 is already in use. Stopping existing process..."
    pkill -f "price-tracker" 2>/dev/null || true
    sleep 2
fi

if check_port 3000; then
    print_warning "Port 3000 is already in use. Stopping existing process..."
    pkill -f "react-scripts" 2>/dev/null || true
    sleep 2
fi

# Install Go dependencies
print_info "Installing Go dependencies..."
go mod tidy
print_status "Go dependencies installed"

# Install UI dependencies
print_info "Installing UI dependencies..."
cd ui
npm install
cd ..
print_status "UI dependencies installed"

# Build the backend
print_info "Building backend application..."
go build -o price-tracker main.go
print_status "Backend built successfully"

# Start the backend
print_info "Starting backend server..."
./price-tracker &
BACKEND_PID=$!
print_status "Backend started with PID: $BACKEND_PID"

# Wait for backend to be ready
if ! wait_for_server 8080; then
    print_error "Backend failed to start properly"
    exit 1
fi

# Start the UI
print_info "Starting React UI..."
cd ui
npm start &
UI_PID=$!
cd ..
print_status "UI started with PID: $UI_PID"

# Wait for UI to be ready
print_info "Waiting for UI to be ready..."
sleep 5

if check_port 3000; then
    print_status "UI is ready on http://localhost:3000"
else
    print_warning "UI may still be starting up..."
fi

echo ""
echo "🎉 Price Tracker Application Started Successfully!"
echo "=================================================="
echo "📊 Backend API: http://localhost:8080"
echo "🌐 React UI:    http://localhost:3000"
echo ""
echo "📋 Available endpoints:"
echo "   - GET  /health          - Health check"
echo "   - GET  /api/candles     - Get candlestick data"
echo "   - WS   /ws/trades       - WebSocket for live trades"
echo ""
echo "🛑 Press Ctrl+C to stop all services"
echo ""

# Keep the script running and monitor processes
while true; do
    # Check if backend is still running
    if ! kill -0 $BACKEND_PID 2>/dev/null; then
        print_error "Backend process died unexpectedly"
        break
    fi
    
    # Check if UI is still running
    if ! kill -0 $UI_PID 2>/dev/null; then
        print_error "UI process died unexpectedly"
        break
    fi
    
    sleep 5
done 