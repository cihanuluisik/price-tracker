import React, { useState } from 'react';
import './App.css';
import LiveTradesPage from './components/LiveTradesPage';
import PriceChartsPage from './components/PriceChartsPage';

function App() {
  const [activePopup, setActivePopup] = useState(null);

  const openPopup = (popupType) => {
    setActivePopup(popupType);
  };

  const closePopup = () => {
    setActivePopup(null);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>🚀 Price Tracker Dashboard</h1>
        <p>Real-time cryptocurrency price tracking and analysis</p>
      </header>

      <main className="App-main">
        <div className="landing-content">
          <h2>Welcome to Price Tracker</h2>
          <p>Select an option to get started:</p>
          
          <div className="button-container">
            <button 
              className="action-button live-trades-btn"
              onClick={() => openPopup('liveTrades')}
            >
              📊 Live Trades
            </button>
            
            <button 
              className="action-button charts-btn"
              onClick={() => openPopup('priceCharts')}
            >
              📈 Price Charts
            </button>
          </div>

          <div className="features">
            <div className="feature-card">
              <h3>🔴 Live Trades</h3>
              <p>View real-time cryptocurrency trades with live WebSocket data</p>
              <ul>
                <li>Real-time price updates</li>
                <li>Trade volume tracking</li>
                <li>Multiple symbol support</li>
                <li>Live timestamp data</li>
              </ul>
            </div>

            <div className="feature-card">
              <h3>📈 Price Charts</h3>
              <p>Interactive candlestick charts with historical data analysis</p>
              <ul>
                <li>Candlestick chart visualization</li>
                <li>Multiple time periods (1m, 2m, 5m, 10m)</li>
                <li>Symbol selection</li>
                <li>Date range filtering</li>
              </ul>
            </div>
          </div>
        </div>
      </main>

      {/* Popup Overlays */}
      {activePopup && (
        <div className="popup-overlay" onClick={closePopup}>
          <div className="popup-content" onClick={(e) => e.stopPropagation()}>
            <button className="close-button" onClick={closePopup}>
              ✕
            </button>
            
            {activePopup === 'liveTrades' && (
              <LiveTradesPage />
            )}
            
            {activePopup === 'priceCharts' && (
              <PriceChartsPage />
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
