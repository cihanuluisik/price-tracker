package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"price-tracker/models"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan *models.Trade
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan *models.Trade, 100),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (s *WebSocketServer) Start() {
	go s.run()
}

func (s *WebSocketServer) run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
			log.Printf("Client connected. Total clients: %d", len(s.clients))

		case client := <-s.unregister:
			s.mutex.Lock()
			delete(s.clients, client)
			s.mutex.Unlock()
			client.Close()
			log.Printf("Client disconnected. Total clients: %d", len(s.clients))

		case trade := <-s.broadcast:
			s.mutex.RLock()
			clients := make([]*websocket.Conn, 0, len(s.clients))
			for client := range s.clients {
				clients = append(clients, client)
			}
			clientCount := len(clients)
			s.mutex.RUnlock()

			// Transform trade data for frontend
			tradeData := map[string]interface{}{
				"symbol":    trade.Symbol,
				"price":     trade.Price,
				"quantity":  trade.Quantity,
				"tradeTime": trade.TradeTime,
				"tradeId":   trade.TradeID,
			}

			tradeJSON, err := json.Marshal(tradeData)
			if err != nil {
				log.Printf("Error marshaling trade data: %v", err)
				continue
			}

			// Broadcast to all connected clients
			successCount := 0
			for _, client := range clients {
				err := client.WriteMessage(websocket.TextMessage, tradeJSON)
				if err != nil {
					log.Printf("Error sending message to client: %v", err)
					s.mutex.Lock()
					delete(s.clients, client)
					s.mutex.Unlock()
					client.Close()
				} else {
					successCount++
				}
			}

			// Log broadcast stats occasionally
			if clientCount > 0 && successCount%50 == 0 {
				log.Printf("Broadcasted trade to %d/%d clients: %s @ $%s", successCount, clientCount, trade.Symbol, trade.Price)
			}
		}
	}
}

func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to WebSocket: %v", err)
		return
	}

	s.register <- conn

	// Handle incoming messages (if any)
	go func() {
		defer func() {
			s.unregister <- conn
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			// Echo message back for testing (optional)
			log.Printf("Received message: %s", message)
		}
	}()
}

func (s *WebSocketServer) BroadcastTrade(trade *models.Trade) {
	s.broadcast <- trade
}

func (s *WebSocketServer) GetClientCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.clients)
}
 