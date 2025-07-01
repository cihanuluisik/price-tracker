package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"price-tracker/models"

	"github.com/gorilla/websocket"
)

type BinanceWebSocket struct {
	conn     *websocket.Conn
	symbols  []string
	handlers []TradeHandler
}

type TradeHandler func(*models.Trade)

func NewBinanceWebSocket(symbols []string) *BinanceWebSocket {
	return &BinanceWebSocket{
		symbols:  symbols,
		handlers: make([]TradeHandler, 0),
	}
}

func (b *BinanceWebSocket) AddHandler(handler TradeHandler) {
	b.handlers = append(b.handlers, handler)
}

func (b *BinanceWebSocket) Connect() error {
	streams := make([]string, len(b.symbols))
	for i, symbol := range b.symbols {
		streams[i] = fmt.Sprintf("%s@trade", strings.ToLower(symbol))
	}

	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s", strings.Join(streams, "/"))

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Binance WebSocket: %v", err)
	}

	b.conn = conn
	log.Printf("Connected to Binance WebSocket for symbols: %v", b.symbols)
	return nil
}

func (b *BinanceWebSocket) Listen() {
	defer b.conn.Close()

	log.Printf("Starting to listen for trade data from Binance...")
	tradeCount := 0

	for {
		_, message, err := b.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		trade, err := b.parseTrade(message)
		if err != nil {
			log.Printf("Error parsing trade: %v", err)
			continue
		}

		tradeCount++
		if tradeCount%100 == 0 { // Log every 100th trade
			log.Printf("Processed %d trades, latest: %s @ $%.2f", tradeCount, trade.Symbol, trade.ParsedPrice)
		}

		for _, handler := range b.handlers {
			handler(trade)
		}
	}
}

func (b *BinanceWebSocket) parseTrade(message []byte) (*models.Trade, error) {
	var trade models.Trade
	if err := json.Unmarshal(message, &trade); err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(trade.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price: %v", err)
	}

	quantity, err := strconv.ParseFloat(trade.Quantity, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quantity: %v", err)
	}

	trade.ParsedPrice = price
	trade.ParsedQuantity = quantity
	trade.ParsedTime = time.Unix(0, trade.TradeTime*int64(time.Millisecond))

	return &trade, nil
}

func (b *BinanceWebSocket) Close() error {
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}
 