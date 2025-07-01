package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"price-tracker/api"
	"price-tracker/config"
	"price-tracker/database"
	"price-tracker/models"
	"price-tracker/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewSQLiteDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize WebSocket server for frontend
	wsServer := websocket.NewWebSocketServer()
	wsServer.Start()

	// Initialize Binance WebSocket client
	binanceWS := websocket.NewBinanceWebSocket(cfg.Symbols)

	// Add handler to save trades and broadcast to frontend
	binanceWS.AddHandler(func(trade *models.Trade) {
		// Save trade to database
		if err := db.SaveTrade(trade); err != nil {
			log.Printf("Failed to save trade: %v", err)
		}

		// Broadcast to frontend
		wsServer.BroadcastTrade(trade)
	})

	// Start Binance WebSocket connection
	go func() {
		if err := binanceWS.Connect(); err != nil {
			log.Printf("Failed to connect to Binance WebSocket: %v", err)
			return
		}
		binanceWS.Listen()
	}()

	// Initialize HTTP handler
	h := api.NewHandler(db)

	// Set up HTTP routes
	http.HandleFunc("/health", h.HealthCheck)
	http.HandleFunc("/api/candles", h.GetCandles)
	http.HandleFunc("/api/realtime-candles", h.GetRealtimeCandles)
	http.HandleFunc("/api/trade-stats", h.GetTradeStats)
	http.HandleFunc("/ws/trades", wsServer.HandleWebSocket)

	// Start HTTP server
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s:%d", cfg.Server.Host, cfg.Server.Port)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	binanceWS.Close()
}
 