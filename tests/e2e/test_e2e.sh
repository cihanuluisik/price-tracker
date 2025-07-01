#!/bin/bash

# E2E Test Runner for Price Tracker Backend Service

set -e

# Change to project root directory
cd "$(dirname "$0")/../.."

echo "🧪 Starting E2E Tests for Price Tracker Backend Service"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up..."
    pkill -f "test-price-tracker" || true
    rm -f test-price-tracker
    rm -f price_tracker.db
    print_status "Cleanup completed"
}

# Set up trap to ensure cleanup on exit
trap cleanup EXIT

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go first."
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    print_error "main.go not found. Please run this script from the project root."
    exit 1
fi

print_status "Go version: $(go version)"

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy
print_status "Dependencies installed"

# Run tests
echo "🚀 Running E2E tests..."
go test -v -timeout=5m ./tests/e2e/e2e_test.go

if [ $? -eq 0 ]; then
    print_status "All E2E tests passed!"
    echo "=================================================="
    echo "🎉 E2E test suite completed successfully!"
else
    print_error "Some E2E tests failed!"
    exit 1
fi 