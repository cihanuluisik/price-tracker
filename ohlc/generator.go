package ohlc

import (
	"math"
	"sync"
	"time"

	"price-tracker/models"
)

type Generator struct {
	periods map[string]map[int]*PeriodData
	mutex   sync.RWMutex
}

type PeriodData struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	OpenTime  time.Time
	CloseTime time.Time
	IsNew     bool
}

func NewGenerator() *Generator {
	return &Generator{
		periods: make(map[string]map[int]*PeriodData),
	}
}

func (g *Generator) ProcessTrade(trade *models.Trade, periods []int) []*models.OHLC {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var ohlcData []*models.OHLC

	for _, period := range periods {
		periodKey := getPeriodKey(trade.Symbol, period)

		if g.periods[periodKey] == nil {
			g.periods[periodKey] = make(map[int]*PeriodData)
		}

		currentData := g.periods[periodKey][period]
		periodStart := getPeriodStart(trade.ParsedTime, period)

		if currentData == nil || !isSamePeriod(currentData.OpenTime, periodStart, period) {
			if currentData != nil && currentData.IsNew {
				ohlcData = append(ohlcData, g.createOHLC(trade.Symbol, period, currentData))
			}

			currentData = &PeriodData{
				Open:     trade.ParsedPrice,
				High:     trade.ParsedPrice,
				Low:      trade.ParsedPrice,
				Close:    trade.ParsedPrice,
				Volume:   trade.ParsedQuantity,
				OpenTime: periodStart,
				IsNew:    true,
			}
		} else {
			currentData.High = math.Max(currentData.High, trade.ParsedPrice)
			currentData.Low = math.Min(currentData.Low, trade.ParsedPrice)
			currentData.Close = trade.ParsedPrice
			currentData.Volume += trade.ParsedQuantity
		}

		currentData.CloseTime = trade.ParsedTime
		g.periods[periodKey][period] = currentData
	}

	return ohlcData
}

func (g *Generator) createOHLC(symbol string, period int, data *PeriodData) *models.OHLC {
	return &models.OHLC{
		Symbol:    symbol,
		Period:    period,
		Open:      data.Open,
		High:      data.High,
		Low:       data.Low,
		Close:     data.Close,
		Volume:    data.Volume,
		OpenTime:  data.OpenTime,
		CloseTime: data.CloseTime,
	}
}

func getPeriodKey(symbol string, period int) string {
	return symbol
}

func getPeriodStart(tradeTime time.Time, periodMinutes int) time.Time {
	minutes := tradeTime.Minute() - (tradeTime.Minute() % periodMinutes)
	return time.Date(tradeTime.Year(), tradeTime.Month(), tradeTime.Day(), tradeTime.Hour(), minutes, 0, 0, tradeTime.Location())
}

func isSamePeriod(time1, time2 time.Time, periodMinutes int) bool {
	period1 := getPeriodStart(time1, periodMinutes)
	period2 := getPeriodStart(time2, periodMinutes)
	return period1.Equal(period2)
}
