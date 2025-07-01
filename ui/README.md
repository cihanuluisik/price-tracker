# Price Tracker UI

A modern React.js frontend for the Price Tracker Backend Service, featuring real-time cryptocurrency trade data and interactive candlestick charts.

## Features

- 🚀 **Modern React UI** - Built with React 18 and modern JavaScript (ES8+)
- 📊 **Live Trades Page** - Real-time trade data with WebSocket connection
- 📈 **Price Charts Page** - Interactive candlestick charts with multiple time periods
- 🎨 **Beautiful Design** - Modern, responsive UI with smooth animations
- 📱 **Mobile Responsive** - Works perfectly on all device sizes
- ⚡ **Real-time Updates** - Live data streaming from Binance WebSocket
- 🔄 **Auto-reconnection** - Automatic WebSocket reconnection on disconnection

## Tech Stack

- **React 18** - Modern React with hooks
- **JavaScript ES8+** - Latest stable JavaScript features
- **Lightweight Charts** - High-performance candlestick charts
- **Axios** - HTTP client for API requests
- **WebSocket API** - Real-time data streaming
- **CSS3** - Modern styling with gradients and animations

## Prerequisites

- Node.js 16+ 
- npm or yarn
- Running Go backend service (on localhost:8080)

## Installation

1. Navigate to the UI directory:
```bash
cd ui
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm start
```

The application will open at `http://localhost:3000`

## Project Structure

```
ui/
├── public/                 # Static files
├── src/
│   ├── components/         # React components
│   │   ├── LiveTradesPage.js      # Live trades component
│   │   ├── LiveTradesPage.css     # Live trades styles
│   │   ├── PriceChartsPage.js     # Price charts component
│   │   └── PriceChartsPage.css    # Price charts styles
│   ├── App.js              # Main application component
│   ├── App.css             # Main application styles
│   ├── index.js            # Application entry point
│   └── index.css           # Global styles
├── package.json            # Dependencies and scripts
└── README.md               # This file
```

## Usage

### Landing Page
- Modern dashboard with feature overview
- Two main action buttons for Live Trades and Price Charts
- Responsive design with feature cards

### Live Trades Page
- **Real-time Trade Table**: Displays live trade data with columns:
  - Symbol (BTC, ETH, etc.)
  - Price (with color-coded changes)
  - Quantity
  - Time
  - Trade ID
- **Connection Status**: Shows WebSocket connection status
- **Trade Summary**: Displays total trades, active symbols, and latest update
- **Auto-reconnection**: Automatically reconnects if connection is lost

### Price Charts Page
- **Interactive Candlestick Charts**: Professional trading charts
- **Symbol Selection**: Choose from BTC, ETH, ADA, DOT, LINK
- **Time Period Selection**: 1m, 2m, 5m, 10m intervals
- **Chart Controls**: Refresh button and real-time updates
- **Chart Information**: Current symbol, period, data points, and last update
- **Responsive Design**: Charts adapt to screen size

## API Integration

### Backend Connection
- **Proxy Configuration**: Automatically proxies API calls to `localhost:8080`
- **WebSocket Connection**: Connects to backend WebSocket for live data
- **REST API**: Fetches historical chart data from backend

### Endpoints Used
- `GET /api/candles` - Fetch candlestick data
- `GET /health` - Health check (via proxy)
- `WebSocket /ws/trades` - Real-time trade data

## Development

### Available Scripts

```bash
# Start development server
npm start

# Build for production
npm run build

# Run tests
npm test

# Eject from Create React App
npm run eject
```

### Development Server
- Runs on `http://localhost:3000`
- Hot reload enabled
- Proxy configured to backend on port 8080

### Building for Production
```bash
npm run build
```
Creates optimized build in `build/` directory

## Configuration

### Backend URL
The UI is configured to connect to the backend at `localhost:8080`. To change this:

1. Update the proxy in `package.json`:
```json
{
  "proxy": "http://your-backend-url:port"
}
```

2. Update WebSocket URL in `LiveTradesPage.js`:
```javascript
const ws = new WebSocket('ws://your-backend-url:port/ws/trades');
```

### Chart Configuration
Chart settings can be modified in `PriceChartsPage.js`:
- Chart height and width
- Color schemes
- Time scale options
- Candlestick styling

## Features in Detail

### Live Trades Features
- **Real-time Updates**: Trades appear instantly as they occur
- **Price Color Coding**: Green for price increases, red for decreases
- **Trade History**: Keeps last 100 trades in memory
- **Connection Monitoring**: Visual indicator of WebSocket status
- **Auto-scroll**: New trades appear at the top
- **Responsive Table**: Adapts to different screen sizes

### Price Charts Features
- **Professional Charts**: Using TradingView's lightweight-charts
- **Multiple Timeframes**: 1, 2, 5, and 10-minute intervals
- **Symbol Switching**: Instant chart updates when symbol changes
- **Interactive Elements**: Hover tooltips, zoom, pan
- **Data Validation**: Error handling for missing or invalid data
- **Loading States**: Visual feedback during data fetching

## Troubleshooting

### Common Issues

1. **Backend Not Running**
   - Ensure Go backend is running on port 8080
   - Check backend logs for errors

2. **WebSocket Connection Failed**
   - Verify backend WebSocket endpoint is available
   - Check browser console for connection errors

3. **Chart Data Not Loading**
   - Ensure backend API is responding
   - Check network tab for API errors
   - Verify symbol and period parameters

4. **Build Errors**
   - Clear node_modules and reinstall: `rm -rf node_modules && npm install`
   - Check Node.js version compatibility

### Debug Mode
Enable debug logging in browser console:
```javascript
// In browser console
localStorage.setItem('debug', 'true');
```

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Performance

- **Lightweight Charts**: High-performance chart rendering
- **Efficient Updates**: Minimal re-renders with React hooks
- **Memory Management**: Automatic cleanup of WebSocket connections
- **Optimized Builds**: Production builds are minified and optimized

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is part of the Price Tracker application suite.

## Support

For issues and questions:
1. Check the troubleshooting section
2. Review browser console for errors
3. Verify backend service is running
4. Check network connectivity
