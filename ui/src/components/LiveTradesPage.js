import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import './LiveTradesPage.css';

// Memoized Trade Row Component to prevent unnecessary re-renders
const TradeRow = React.memo(({ trade, previousTrade, index }) => {
  const formatTime = useCallback((timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString();
  }, []);

  const formatPrice = useCallback((price) => {
    return parseFloat(price).toFixed(2);
  }, []);

  const formatQuantity = useCallback((quantity) => {
    return parseFloat(quantity).toFixed(4);
  }, []);

  const getPriceColor = useCallback((price, previousPrice) => {
    if (!previousPrice) return 'neutral';
    return parseFloat(price) > parseFloat(previousPrice) ? 'positive' : 'negative';
  }, []);

  const priceColor = getPriceColor(trade.price, previousTrade?.price);
  
  return (
    <tr className="trade-row">
      <td className="symbol-cell">
        <span className="symbol-badge">{trade.symbol}</span>
      </td>
      <td className={`price-cell ${priceColor}`}>
        ${formatPrice(trade.price)}
      </td>
      <td className="quantity-cell">
        {formatQuantity(trade.quantity)}
      </td>
      <td className="time-cell">
        {formatTime(trade.tradeTime)}
      </td>
      <td className="trade-id-cell">
        {trade.tradeId}
      </td>
    </tr>
  );
});

// Memoized Trades Table Component
const TradesTable = React.memo(({ trades, isConnected }) => {
  // Sort trades by symbol to prevent table shaking
  const sortedTrades = useMemo(() => {
    return [...trades].sort((a, b) => {
      // First sort by symbol alphabetically
      const symbolComparison = a.symbol.localeCompare(b.symbol);
      if (symbolComparison !== 0) {
        return symbolComparison;
      }
      // If symbols are the same, sort by trade time (newest first)
      return new Date(b.tradeTime) - new Date(a.tradeTime);
    });
  }, [trades]);

  return (
    <div className="trades-table-container">
      <table className="trades-table">
        <thead>
          <tr>
            <th>Symbol</th>
            <th>Price</th>
            <th>Quantity</th>
            <th>Time</th>
            <th>Trade ID</th>
          </tr>
        </thead>
        <tbody>
          {sortedTrades.length === 0 ? (
            <tr>
              <td colSpan="5" className="no-data">
                {isConnected ? 'Waiting for trade data...' : 'Connecting to server...'}
              </td>
            </tr>
          ) : (
            sortedTrades.map((trade, index) => {
              const previousTrade = sortedTrades[index + 1];
              // Use a stable key based on trade ID and received timestamp
              const stableKey = `${trade.symbol}-${trade.tradeId}-${trade.receivedAt || trade.tradeTime}`;
              
              return (
                <TradeRow
                  key={stableKey}
                  trade={trade}
                  previousTrade={previousTrade}
                  index={index}
                />
              );
            })
          )}
        </tbody>
      </table>
    </div>
  );
});

// Memoized Summary Component
const TradesSummary = React.memo(({ trades }) => {
  const formatTime = useCallback((timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString();
  }, []);

  const summaryData = useMemo(() => ({
    totalTrades: trades.length,
    activeSymbols: new Set(trades.map(trade => trade.symbol)).size,
    latestUpdate: trades.length > 0 ? formatTime(trades[0].tradeTime) : 'N/A'
  }), [trades, formatTime]);

  return (
    <div className="trades-summary">
      <div className="summary-card">
        <h3>Total Trades</h3>
        <p className="summary-value">{summaryData.totalTrades}</p>
      </div>
      
      <div className="summary-card">
        <h3>Active Symbols</h3>
        <p className="summary-value">{summaryData.activeSymbols}</p>
      </div>
      
      <div className="summary-card">
        <h3>Latest Update</h3>
        <p className="summary-value">{summaryData.latestUpdate}</p>
      </div>
    </div>
  );
});

const LiveTradesPage = () => {
  const [trades, setTrades] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState('Connecting...');
  const wsRef = useRef(null);
  const isConnectedRef = useRef(false);

  const connectWebSocket = useCallback(() => {
    // Don't connect if already connected
    if (isConnectedRef.current) {
      return;
    }

    try {
      // Connect to the backend WebSocket for trade data
      const ws = new WebSocket('ws://localhost:8080/ws/trades');
      wsRef.current = ws;

      ws.onopen = () => {
        isConnectedRef.current = true;
        setIsConnected(true);
        setConnectionStatus('Connected');
        console.log('WebSocket connected');
      };

      ws.onmessage = (event) => {
        try {
          const trade = JSON.parse(event.data);
          // Add timestamp for stable keys
          trade.receivedAt = Date.now();
          setTrades(prevTrades => {
            // Use functional update to ensure we're working with the latest state
            const newTrades = [trade, ...prevTrades.slice(0, 99)]; // Keep last 100 trades
            return newTrades;
          });
        } catch (error) {
          console.error('Error parsing trade data:', error);
        }
      };

      ws.onclose = () => {
        isConnectedRef.current = false;
        setIsConnected(false);
        setConnectionStatus('Disconnected');
        console.log('WebSocket disconnected');
        
        // Attempt to reconnect after 3 seconds
        setTimeout(() => {
          if (!isConnectedRef.current) {
            connectWebSocket();
          }
        }, 3000);
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus('Connection Error');
      };

    } catch (error) {
      console.error('Failed to connect to WebSocket:', error);
      setConnectionStatus('Connection Failed');
    }
  }, []);

  useEffect(() => {
    connectWebSocket();
    return () => {
      isConnectedRef.current = false;
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connectWebSocket]);

  return (
    <div className="live-trades-page">
      <div className="page-header">
        <h2>📊 Live Trades</h2>
        <div className="connection-status">
          <span className={`status-indicator ${isConnected ? 'connected' : 'disconnected'}`}>
            ●
          </span>
          <span className="status-text">{connectionStatus}</span>
        </div>
      </div>

      <div className="trades-container">
        <TradesTable trades={trades} isConnected={isConnected} />
        <TradesSummary trades={trades} />
      </div>

      <div className="page-footer">
        <p>Real-time cryptocurrency trade data from Binance WebSocket streams</p>
        <p>Data updates automatically every time a new trade occurs</p>
      </div>
    </div>
  );
};

export default LiveTradesPage; 