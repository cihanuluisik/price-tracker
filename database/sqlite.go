package database

import (
	"database/sql"
	"time"

	"price-tracker/models"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS trades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		trade_id INTEGER NOT NULL,
		price REAL NOT NULL,
		quantity REAL NOT NULL,
		buyer_order_id INTEGER,
		seller_order_id INTEGER,
		trade_time DATETIME NOT NULL,
		is_buyer_maker BOOLEAN,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(symbol, trade_id)
	);
	CREATE INDEX IF NOT EXISTS idx_trades_symbol_time ON trades(symbol, trade_time);
	CREATE INDEX IF NOT EXISTS idx_trades_trade_id ON trades(trade_id);
	
	CREATE TABLE IF NOT EXISTS ohlc (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		period INTEGER NOT NULL,
		open REAL NOT NULL,
		high REAL NOT NULL,
		low REAL NOT NULL,
		close REAL NOT NULL,
		volume REAL NOT NULL,
		open_time DATETIME NOT NULL,
		close_time DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(symbol, period, open_time)
	);
	CREATE INDEX IF NOT EXISTS idx_symbol_period_time ON ohlc(symbol, period, open_time);
	`

	_, err := db.Exec(query)
	return err
}

func (s *SQLiteDB) SaveOHLC(ohlc *models.OHLC) error {
	query := `
	INSERT OR REPLACE INTO ohlc (symbol, period, open, high, low, close, volume, open_time, close_time)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, ohlc.Symbol, ohlc.Period, ohlc.Open, ohlc.High, ohlc.Low, ohlc.Close, ohlc.Volume, ohlc.OpenTime, ohlc.CloseTime)
	return err
}

func (s *SQLiteDB) GetLastHourCandles(symbol string, period int) ([]models.CandleStick, error) {
	query := `
	SELECT symbol, period, open, high, low, close, volume, open_time, close_time
	FROM ohlc
	WHERE symbol = ? AND period = ?
	AND close_time >= datetime('now', '-1 hour', 'localtime')
	ORDER BY open_time DESC
	`

	rows, err := s.db.Query(query, symbol, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []models.CandleStick
	for rows.Next() {
		var candle models.CandleStick
		var openTimeStr, closeTimeStr string

		err := rows.Scan(&candle.Symbol, &candle.Period, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume, &openTimeStr, &closeTimeStr)
		if err != nil {
			return nil, err
		}

		// Try parsing ISO format first (with timezone)
		candle.OpenTime, err = time.Parse(time.RFC3339, openTimeStr)
		if err != nil {
			// Try without timezone
			candle.OpenTime, err = time.Parse("2006-01-02T15:04:05", openTimeStr)
			if err != nil {
				// Try legacy format
				candle.OpenTime, err = time.Parse("2006-01-02 15:04:05", openTimeStr)
				if err != nil {
					return nil, err
				}
			}
		}

		candle.CloseTime, err = time.Parse(time.RFC3339, closeTimeStr)
		if err != nil {
			// Try without timezone
			candle.CloseTime, err = time.Parse("2006-01-02T15:04:05", closeTimeStr)
			if err != nil {
				// Try legacy format
				candle.CloseTime, err = time.Parse("2006-01-02 15:04:05", closeTimeStr)
				if err != nil {
					return nil, err
				}
			}
		}

		candles = append(candles, candle)
	}

	return candles, nil
}

func (s *SQLiteDB) GetLatestTradeTime(symbol string) (time.Time, error) {
	query := `
	SELECT MAX(close_time)
	FROM ohlc
	WHERE symbol = ?
	`

	var latestTimeStr sql.NullString
	err := s.db.QueryRow(query, symbol).Scan(&latestTimeStr)
	if err != nil {
		return time.Time{}, err
	}

	if !latestTimeStr.Valid {
		return time.Time{}, nil
	}

	return time.Parse("2006-01-02 15:04:05", latestTimeStr.String)
}

// SaveTrade saves a trade to the database
func (s *SQLiteDB) SaveTrade(trade *models.Trade) error {
	query := `
	INSERT OR REPLACE INTO trades (symbol, trade_id, price, quantity, buyer_order_id, seller_order_id, trade_time, is_buyer_maker)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		trade.Symbol,
		trade.TradeID,
		trade.ParsedPrice,
		trade.ParsedQuantity,
		trade.BuyerOrderID,
		trade.SellerOrderID,
		trade.ParsedTime,
		trade.IsBuyerMaker,
	)
	return err
}

// GenerateCandlesFromTrades generates candlestick data from stored trades for a given period
func (s *SQLiteDB) GenerateCandlesFromTrades(symbol string, period int, hours int) ([]models.CandleStick, error) {
	// Calculate the time range
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	// Get trades for the specified time range
	query := `
	SELECT trade_time, price, quantity
	FROM trades
	WHERE symbol = ? AND trade_time >= ? AND trade_time <= ?
	ORDER BY trade_time ASC
	`

	rows, err := s.db.Query(query, symbol, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []struct {
		TradeTime time.Time
		Price     float64
		Quantity  float64
	}

	for rows.Next() {
		var trade struct {
			TradeTime time.Time
			Price     float64
			Quantity  float64
		}
		var tradeTimeStr string

		err := rows.Scan(&tradeTimeStr, &trade.Price, &trade.Quantity)
		if err != nil {
			return nil, err
		}

		// Parse time
		trade.TradeTime, err = time.Parse("2006-01-02 15:04:05", tradeTimeStr)
		if err != nil {
			// Try RFC3339 format
			trade.TradeTime, err = time.Parse(time.RFC3339, tradeTimeStr)
			if err != nil {
				return nil, err
			}
		}

		trades = append(trades, trade)
	}

	// Generate candlesticks
	return s.generateCandlesticks(trades, symbol, period), nil
}

// generateCandlesticks creates candlestick data from a slice of trades
func (s *SQLiteDB) generateCandlesticks(trades []struct {
	TradeTime time.Time
	Price     float64
	Quantity  float64
}, symbol string, period int) []models.CandleStick {
	if len(trades) == 0 {
		return []models.CandleStick{}
	}

	var candles []models.CandleStick
	var periodDuration time.Duration

	// Period is now always in seconds
	periodDuration = time.Duration(period) * time.Second

	// Group trades by time period
	currentPeriod := trades[0].TradeTime.Truncate(periodDuration)
	var currentCandle *models.CandleStick

	for _, trade := range trades {
		tradePeriod := trade.TradeTime.Truncate(periodDuration)

		// If we're in a new period, save the current candle and start a new one
		if tradePeriod.After(currentPeriod) {
			if currentCandle != nil {
				candles = append(candles, *currentCandle)
			}
			currentPeriod = tradePeriod
			currentCandle = &models.CandleStick{
				Symbol:    symbol,
				Period:    period,
				Open:      trade.Price,
				High:      trade.Price,
				Low:       trade.Price,
				Close:     trade.Price,
				Volume:    trade.Quantity,
				OpenTime:  tradePeriod,
				CloseTime: tradePeriod.Add(periodDuration),
			}
		} else {
			// Update current candle
			if currentCandle == nil {
				currentCandle = &models.CandleStick{
					Symbol:    symbol,
					Period:    period,
					Open:      trade.Price,
					High:      trade.Price,
					Low:       trade.Price,
					Close:     trade.Price,
					Volume:    trade.Quantity,
					OpenTime:  tradePeriod,
					CloseTime: tradePeriod.Add(periodDuration),
				}
			} else {
				if trade.Price > currentCandle.High {
					currentCandle.High = trade.Price
				}
				if trade.Price < currentCandle.Low {
					currentCandle.Low = trade.Price
				}
				currentCandle.Close = trade.Price
				currentCandle.Volume += trade.Quantity
			}
		}
	}

	// Add the last candle if it exists
	if currentCandle != nil {
		candles = append(candles, *currentCandle)
	}

	// Reverse the candles to get chronological order (oldest first)
	for i, j := 0, len(candles)-1; i < j; i, j = i+1, j-1 {
		candles[i], candles[j] = candles[j], candles[i]
	}

	return candles
}

// GetTradeCount returns the number of trades for a symbol
func (s *SQLiteDB) GetTradeCount(symbol string) (int, error) {
	query := `SELECT COUNT(*) FROM trades WHERE symbol = ?`
	var count int
	err := s.db.QueryRow(query, symbol).Scan(&count)
	return count, err
}

// GetLatestTrade returns the most recent trade for a symbol
func (s *SQLiteDB) GetLatestTrade(symbol string) (*models.Trade, error) {
	query := `
	SELECT symbol, trade_id, price, quantity, buyer_order_id, seller_order_id, trade_time, is_buyer_maker
	FROM trades
	WHERE symbol = ?
	ORDER BY trade_time DESC
	LIMIT 1
	`

	var trade models.Trade
	var tradeTimeStr string

	err := s.db.QueryRow(query, symbol).Scan(
		&trade.Symbol,
		&trade.TradeID,
		&trade.ParsedPrice,
		&trade.ParsedQuantity,
		&trade.BuyerOrderID,
		&trade.SellerOrderID,
		&tradeTimeStr,
		&trade.IsBuyerMaker,
	)

	if err != nil {
		return nil, err
	}

	// Parse time
	trade.ParsedTime, err = time.Parse("2006-01-02 15:04:05", tradeTimeStr)
	if err != nil {
		// Try RFC3339 format
		trade.ParsedTime, err = time.Parse(time.RFC3339, tradeTimeStr)
		if err != nil {
			return nil, err
		}
	}

	return &trade, nil
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}
 