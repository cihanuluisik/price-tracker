package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"price-tracker/config"
	"price-tracker/database"
)

type Handler struct {
	db *database.SQLiteDB
}

func NewHandler(db *database.SQLiteDB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetCandles(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	periodStr := r.URL.Query().Get("period")
	if periodStr == "" {
		http.Error(w, "period parameter is required", http.StatusBadRequest)
		return
	}

	// Parse period string (e.g., "10s", "1m", "5m")
	periodInfo, err := config.GetPeriodInfo(periodStr)
	if err != nil {
		http.Error(w, "invalid period parameter: "+err.Error(), http.StatusBadRequest)
		return
	}
	period := periodInfo.Value // seconds

	// Get time range parameters
	hoursStr := r.URL.Query().Get("hours")
	hours := 1 // default to 1 hour
	if hoursStr != "" {
		if hours, err = strconv.Atoi(hoursStr); err != nil {
			http.Error(w, "invalid hours parameter", http.StatusBadRequest)
			return
		}
	}

	// Adjust time range for shorter periods (seconds)
	if period < 60 {
		if hours > 6 {
			hours = 6
		}
	}

	candles, err := h.db.GenerateCandlesFromTrades(symbol, period, hours)
	if err != nil {
		http.Error(w, "failed to generate candles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tradeCount, _ := h.db.GetTradeCount(symbol)
	latestTrade, _ := h.db.GetLatestTrade(symbol)

	response := map[string]interface{}{
		"symbol":     symbol,
		"period":     periodInfo.Original,
		"candles":    candles,
		"tradeCount": tradeCount,
		"generated":  "dynamic",
		"hours":      hours,
	}

	if latestTrade != nil {
		response["latestTrade"] = map[string]interface{}{
			"price":     latestTrade.ParsedPrice,
			"quantity":  latestTrade.ParsedQuantity,
			"tradeTime": latestTrade.ParsedTime,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRealtimeCandles returns the most recent candlestick data with real-time updates
func (h *Handler) GetRealtimeCandles(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	periodStr := r.URL.Query().Get("period")
	if periodStr == "" {
		http.Error(w, "period parameter is required", http.StatusBadRequest)
		return
	}

	periodInfo, err := config.GetPeriodInfo(periodStr)
	if err != nil {
		http.Error(w, "invalid period parameter: "+err.Error(), http.StatusBadRequest)
		return
	}
	period := periodInfo.Value

	candles, err := h.db.GenerateCandlesFromTrades(symbol, period, 1) // 1 hour
	if err != nil {
		http.Error(w, "failed to generate candles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	latestTrade, _ := h.db.GetLatestTrade(symbol)
	tradeCount, _ := h.db.GetTradeCount(symbol)

	response := map[string]interface{}{
		"symbol":     symbol,
		"period":     periodInfo.Original,
		"candles":    candles,
		"tradeCount": tradeCount,
		"generated":  "realtime",
		"timestamp":  time.Now().Unix(),
	}

	if latestTrade != nil {
		response["latestTrade"] = map[string]interface{}{
			"price":     latestTrade.ParsedPrice,
			"quantity":  latestTrade.ParsedQuantity,
			"tradeTime": latestTrade.ParsedTime,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// GetTradeStats returns statistics about stored trades
func (h *Handler) GetTradeStats(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	tradeCount, err := h.db.GetTradeCount(symbol)
	if err != nil {
		http.Error(w, "failed to get trade count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	latestTrade, err := h.db.GetLatestTrade(symbol)
	if err != nil && err.Error() != "sql: no rows in result set" {
		http.Error(w, "failed to get latest trade: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"symbol":     symbol,
		"tradeCount": tradeCount,
	}

	if latestTrade != nil {
		response["latestTrade"] = map[string]interface{}{
			"price":     latestTrade.ParsedPrice,
			"quantity":  latestTrade.ParsedQuantity,
			"tradeTime": latestTrade.ParsedTime,
			"tradeId":   latestTrade.TradeID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
