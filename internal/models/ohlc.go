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
	Symbol   string `json:"symbol"`
	Period   int    `json:"period"`
	Candles  []OHLC `json:"candles"`
	Count    int    `json:"count"`
	FromTime string `json:"from_time"`
	ToTime   string `json:"to_time"`
}
