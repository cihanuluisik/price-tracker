package models

import "time"

type Trade struct {
	Symbol         string `json:"s"`
	TradeID        int64  `json:"t"`
	Price          string `json:"p"`
	Quantity       string `json:"q"`
	BuyerOrderID   int64  `json:"b"`
	SellerOrderID  int64  `json:"a"`
	TradeTime      int64  `json:"T"`
	IsBuyerMaker   bool   `json:"m"`
	Ignore         bool   `json:"M"`
	ParsedPrice    float64
	ParsedQuantity float64
	ParsedTime     time.Time
}

type OHLC struct {
	Symbol    string    `json:"symbol"`
	Period    int       `json:"period"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}

type CandleStick struct {
	Symbol    string    `json:"symbol"`
	Period    int       `json:"period"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}
