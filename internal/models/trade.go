package models

import (
	"strconv"
	"time"
)

type Trade struct {
	Symbol         string `json:"s"`
	ID             int64  `json:"t"`
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

func (t *Trade) Parse() error {
	price, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		return err
	}
	t.ParsedPrice = price

	quantity, err := strconv.ParseFloat(t.Quantity, 64)
	if err != nil {
		return err
	}
	t.ParsedQuantity = quantity

	// Convert milliseconds to time.Time
	t.ParsedTime = time.Unix(0, t.TradeTime*int64(time.Millisecond))
	return nil
}
