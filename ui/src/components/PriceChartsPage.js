import React, { useState, useEffect, useRef, useCallback } from 'react';
import { createChart } from 'lightweight-charts';
import axios from 'axios';
import './PriceChartsPage.css';

const PriceChartsPage = () => {
  const [selectedSymbol, setSelectedSymbol] = useState('BTCUSDT');
  const [selectedPeriod, setSelectedPeriod] = useState('10s');
  const [chartData, setChartData] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const chartContainerRef = useRef(null);
  const chartRef = useRef(null);

  const symbols = ['BTCUSDT', 'ETHUSDT', 'ADAUSDT', 'DOTUSDT', 'LINKUSDT'];
  const periods = [
    { value: '10s', label: '10 Seconds' },
    { value: '30s', label: '30 Seconds' },
    { value: '1m', label: '1 Minute' },
    { value: '2m', label: '2 Minutes' },
    { value: '5m', label: '5 Minutes' }
  ];

  const fetchChartData = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await axios.get(`/api/candles`, {
        params: {
          symbol: selectedSymbol,
          period: selectedPeriod
        }
      });

      const candles = response.data.candles || [];
      
      // Transform data for lightweight-charts
      const transformedData = candles.map(candle => ({
        time: Math.floor(new Date(candle.open_time).getTime() / 1000),
        open: parseFloat(candle.open),
        high: parseFloat(candle.high),
        low: parseFloat(candle.low),
        close: parseFloat(candle.close)
      })).reverse(); // Reverse to show oldest to newest

      setChartData(transformedData);

      // Update chart
      if (chartRef.current && chartRef.current.candlestickSeries) {
        try {
          chartRef.current.candlestickSeries.setData(transformedData);
        } catch (error) {
          console.log('Chart disposed during data update, ignoring...');
        }
      }

    } catch (err) {
      console.error('Error fetching chart data:', err);
      setError('Failed to load chart data. Please try again.');
    } finally {
      setIsLoading(false);
    }
  }, [selectedSymbol, selectedPeriod]);

  useEffect(() => {
    if (chartContainerRef.current) {
      const cleanup = initializeChart();
      return cleanup;
    }
  }, []);

  useEffect(() => {
    fetchChartData();
  }, [fetchChartData]);

  const initializeChart = () => {
    // Only remove if chart exists and is not already disposed
    if (chartRef.current && chartRef.current.remove) {
      try {
        chartRef.current.remove();
      } catch (error) {
        // Chart might already be disposed, ignore the error
        console.log('Chart already disposed, continuing...');
      }
    }

    // Clear the reference
    chartRef.current = null;

    // Check if container exists
    if (!chartContainerRef.current) {
      return () => {};
    }

    const chart = createChart(chartContainerRef.current, {
      width: chartContainerRef.current.clientWidth,
      height: 500,
      layout: {
        background: { color: '#ffffff' },
        textColor: '#333',
      },
      grid: {
        vertLines: { color: '#f0f0f0' },
        horzLines: { color: '#f0f0f0' },
      },
      crosshair: {
        mode: 1,
      },
      rightPriceScale: {
        borderColor: '#f0f0f0',
      },
      timeScale: {
        borderColor: '#f0f0f0',
        timeVisible: true,
        secondsVisible: false,
      },
    });

    const candlestickSeries = chart.addCandlestickSeries({
      upColor: '#26a69a',
      downColor: '#ef5350',
      borderVisible: false,
      wickUpColor: '#26a69a',
      wickDownColor: '#ef5350',
    });

    chartRef.current = chart;
    chartRef.current.candlestickSeries = candlestickSeries;

    const handleResize = () => {
      if (chartRef.current && chartRef.current.applyOptions) {
        try {
          chartRef.current.applyOptions({
            width: chartContainerRef.current.clientWidth,
          });
        } catch (error) {
          // Chart might be disposed, ignore resize errors
          console.log('Chart disposed during resize, ignoring...');
        }
      }
    };

    window.addEventListener('resize', handleResize);
    
    // Return cleanup function
    return () => {
      window.removeEventListener('resize', handleResize);
      if (chartRef.current && chartRef.current.remove) {
        try {
          chartRef.current.remove();
        } catch (error) {
          // Ignore disposal errors
        }
      }
      chartRef.current = null;
    };
  };

  const formatSymbol = (symbol) => {
    return symbol.replace('USDT', '');
  };

  const getPeriodLabel = (period) => {
    const periodObj = periods.find(p => p.value === period);
    return periodObj ? periodObj.label : period;
  };

  return (
    <div className="price-charts-page">
      <div className="page-header">
        <h2>📈 Price Charts</h2>
        <div className="chart-controls">
          <div className="control-group">
            <label htmlFor="symbol-select">Symbol:</label>
            <select
              id="symbol-select"
              value={selectedSymbol}
              onChange={(e) => setSelectedSymbol(e.target.value)}
              className="control-select"
            >
              {symbols.map(symbol => (
                <option key={symbol} value={symbol}>
                  {formatSymbol(symbol)}
                </option>
              ))}
            </select>
          </div>

          <div className="control-group">
            <label htmlFor="period-select">Period:</label>
            <select
              id="period-select"
              value={selectedPeriod}
              onChange={(e) => setSelectedPeriod(e.target.value)}
              className="control-select"
            >
              {periods.map(period => (
                <option key={period.value} value={period.value}>
                  {period.label}
                </option>
              ))}
            </select>
          </div>

          <button
            onClick={fetchChartData}
            disabled={isLoading}
            className="refresh-button"
          >
            {isLoading ? '🔄 Loading...' : '🔄 Refresh'}
          </button>
        </div>
      </div>

      <div className="chart-info">
        <div className="info-card">
          <h3>Current Symbol</h3>
          <p className="info-value">{formatSymbol(selectedSymbol)}</p>
        </div>
        
        <div className="info-card">
          <h3>Time Period</h3>
          <p className="info-value">{getPeriodLabel(selectedPeriod)}</p>
        </div>
        
        <div className="info-card">
          <h3>Data Points</h3>
          <p className="info-value">{chartData.length}</p>
        </div>
        
        <div className="info-card">
          <h3>Last Update</h3>
          <p className="info-value">
            {chartData.length > 0 
              ? new Date(chartData[chartData.length - 1].time * 1000).toLocaleString()
              : 'N/A'
            }
          </p>
        </div>
      </div>

      <div className="chart-container">
        {error && (
          <div className="error-message">
            <p>❌ {error}</p>
            <button onClick={fetchChartData} className="retry-button">
              Try Again
            </button>
          </div>
        )}

        {isLoading && (
          <div className="loading-message">
            <div className="loading-spinner"></div>
            <p>Loading chart data...</p>
          </div>
        )}

        <div 
          ref={chartContainerRef} 
          className="chart-wrapper"
          style={{ display: isLoading || error ? 'none' : 'block' }}
        />

        {!isLoading && !error && chartData.length === 0 && (
          <div className="no-data-message">
            <p>📊 No chart data available</p>
            <p>Try selecting a different symbol or period</p>
          </div>
        )}
      </div>

      <div className="chart-footer">
        <div className="legend">
          <div className="legend-item">
            <span className="legend-color positive"></span>
            <span>Bullish (Green)</span>
          </div>
          <div className="legend-item">
            <span className="legend-color negative"></span>
            <span>Bearish (Red)</span>
          </div>
        </div>
        
        <div className="data-source">
          <p>📊 Data from Binance WebSocket streams via Go backend</p>
          <p>🕐 Last hour of OHLC data for {formatSymbol(selectedSymbol)} ({getPeriodLabel(selectedPeriod)})</p>
        </div>
      </div>
    </div>
  );
};

export default PriceChartsPage; 