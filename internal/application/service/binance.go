package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type BinanceSvc interface {
	GetPing() (any, error)
	GetServerTime() (any, error)
	GetExchangeInfo() (any, error)
	GetTickerPrice(symbol string) (any, error)
	GetAllTickerPrices() (any, error)
	GetBookTicker(symbol string) (any, error)
	GetDepth(symbol string, limit int) (any, error)
	GetRecentTrades(symbol string, limit int) (any, error)
	GetKlines(symbol, interval string, limit int) (any, error)
	GetHistoricalTrades(symbol string, limit int, fromId *int64) (any, error)
	GetAggregateTrades(symbol string, fromId, startTime, endTime *int64, limit int) (any, error)
	GetAvgPrice(symbol string) (any, error)
	GetTicker24Hr(symbol string) (any, error)
	GetAllBookTickers() (any, error)
}

type binanceSvc struct {
	baseURL       string
	localCacheSvc LocalCacheSvc
	cacheTTL      time.Duration
	cacheDelay    time.Duration
	lock          sync.RWMutex
}

func NewBinanceSvc() BinanceSvc {
	return &binanceSvc{
		baseURL:       "https://api.binance.com",
		localCacheSvc: NewLocalCacheSvc(),
		cacheTTL:      1 * time.Minute,
		cacheDelay:    500 * time.Millisecond,
	}
}

// fetchData makes an HTTP GET request to the given API URL with parameters.
func (s *binanceSvc) fetchData(apiURL string, params map[string]string) (any, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching data from %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code %d from %s, response: %s", resp.StatusCode, u.String(), resp.Status)
	}

	var response any
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response from %s: %w", u.String(), err)
	}
	return response, nil
}

// fetchAndCache fetches data from the API and stores it in the local cache.
func (s *binanceSvc) fetchAndCache(key, delayKey, apiURL string, params map[string]string) (any, error) {
	data, err := s.fetchData(apiURL, params)
	if err != nil {
		return nil, err
	}

	s.localCacheSvc.Set(key, data, s.cacheTTL)
	s.localCacheSvc.Set(delayKey, true, s.cacheDelay)
	return data, nil
}

// refreshCache asynchronously refreshes the cache for a given key if the delay period has passed.
func (s *binanceSvc) refreshCache(key, delayKey, apiURL string, params map[string]string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, delayExists := s.localCacheSvc.Get(delayKey); delayExists {
		return
	}
	s.localCacheSvc.Set(delayKey, true, s.cacheDelay)

	data, err := s.fetchData(apiURL, params)
	if err != nil {
		log.Printf("Failed to refresh spot cache for %s: %v", key, err)
		s.localCacheSvc.Del(delayKey)
		return
	}

	s.localCacheSvc.Set(key, data, s.cacheTTL)
}

// getWithCache retrieves data from cache or fetches it from the API, caching the result.
func (s *binanceSvc) getWithCache(cacheName, keySuffix, apiURL string, params map[string]string) (any, error) {
	key := fmt.Sprintf("spot_%s:%s", cacheName, keySuffix)
	delayKey := fmt.Sprintf("spot_%s:%s:delay", cacheName, keySuffix)

	if cachedData, found := s.localCacheSvc.Get(key); found {
		go s.refreshCache(key, delayKey, apiURL, params)
		return cachedData, nil
	}

	return s.fetchAndCache(key, delayKey, apiURL, params)
}

// General Endpoints (Spot)

// GetPing tests connectivity to the Rest API.
func (s *binanceSvc) GetPing() (any, error) {
	return map[string]any{
		"serverTime": time.Now().UnixMilli(),
		"message":    "success",
	}, nil
}

// GetServerTime tests connectivity to the Rest API and get the current server time.
func (s *binanceSvc) GetServerTime() (any, error) {
	return map[string]any{
		"serverTime": time.Now().UnixMilli(),
	}, nil
}

// GetExchangeInfo current exchange trading rules and symbol information.
func (s *binanceSvc) GetExchangeInfo() (any, error) {
	return s.getWithCache("exchangeinfo", "global", fmt.Sprintf("%v/api/v3/exchangeInfo", s.baseURL), nil)
}

// Market Data Endpoints (Spot)

// GetTickerPrice returns the latest price for a symbol or all symbols.
func (s *binanceSvc) GetTickerPrice(symbol string) (any, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("tickerprice", symbol, fmt.Sprintf("%v/api/v3/ticker/price", s.baseURL), params)
}

// GetAllTickerPrices returns the latest price for all symbols.
func (s *binanceSvc) GetAllTickerPrices() (any, error) {
	return s.getWithCache("alltickerprices", "global", fmt.Sprintf("%v/api/v3/ticker/price", s.baseURL), nil)
}

// GetBookTicker returns the best price/qty on the order book for a symbol.
func (s *binanceSvc) GetBookTicker(symbol string) (any, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("bookticker", symbol, s.baseURL+"/api/v3/ticker/bookTicker", params)
}

// GetDepth returns the order book for a symbol.
func (s *binanceSvc) GetDepth(symbol string, limit int) (any, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}
	return s.getWithCache("depth", fmt.Sprintf("%s-%d", symbol, limit), s.baseURL+"/api/v3/depth", params)
}

// GetRecentTrades Get recent trades.
func (s *binanceSvc) GetRecentTrades(symbol string, limit int) (any, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}
	return s.getWithCache("recenttrades", fmt.Sprintf("%s-%d", symbol, limit), s.baseURL+"/api/v3/trades", params)
}

// GetKlines returns candlestick data for a symbol.
func (s *binanceSvc) GetKlines(symbol, interval string, limit int) (any, error) {
	params := map[string]string{
		"symbol":   symbol,
		"interval": interval,
		"limit":    fmt.Sprintf("%d", limit),
	}
	return s.getWithCache("klines", fmt.Sprintf("%s-%s-%d", symbol, interval, limit), s.baseURL+"/api/v3/klines", params)
}

// GetHistoricalTrades Get compressed, aggregate trades.
func (s *binanceSvc) GetHistoricalTrades(symbol string, limit int, fromId *int64) (any, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}
	if fromId != nil {
		params["fromId"] = fmt.Sprintf("%d", *fromId)
	}
	keySuffix := fmt.Sprintf("%s-%d", symbol, limit)
	if fromId != nil {
		keySuffix += fmt.Sprintf("-%d", *fromId)
	}
	return s.getWithCache("historicaltrades", keySuffix, s.baseURL+"/api/v3/historicalTrades", params)
}

// GetAggregateTrades Get compressed, aggregate trades.
func (s *binanceSvc) GetAggregateTrades(symbol string, fromId, startTime, endTime *int64, limit int) (any, error) {
	params := map[string]string{
		"symbol": symbol,
	}
	if fromId != nil {
		params["fromId"] = fmt.Sprintf("%d", *fromId)
	}
	if startTime != nil {
		params["startTime"] = fmt.Sprintf("%d", *startTime)
	}
	if endTime != nil {
		params["endTime"] = fmt.Sprintf("%d", *endTime)
	}
	params["limit"] = fmt.Sprintf("%d", limit)

	keySuffix := fmt.Sprintf("%s-%d", symbol, limit)
	if fromId != nil {
		keySuffix += fmt.Sprintf("-f%d", *fromId)
	}
	if startTime != nil {
		keySuffix += fmt.Sprintf("-s%d", *startTime)
	}
	if endTime != nil {
		keySuffix += fmt.Sprintf("-e%d", *endTime)
	}
	return s.getWithCache("aggregatetrades", keySuffix, s.baseURL+"/api/v3/aggTrades", params)
}

// GetAvgPrice Current average price for a symbol.
func (s *binanceSvc) GetAvgPrice(symbol string) (any, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("avgprice", symbol, s.baseURL+"/api/v3/avgPrice", params)
}

// GetTicker24Hr 24hr Ticker Price Change Statistics.
func (s *binanceSvc) GetTicker24Hr(symbol string) (any, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("ticker24hr", symbol, s.baseURL+"/api/v3/ticker/24hr", params)
}

// GetAllBookTickers returns the best price/qty on the order book for all symbols.
func (s *binanceSvc) GetAllBookTickers() (any, error) {
	return s.getWithCache("allbooktickers", "global", s.baseURL+"/api/v3/ticker/bookTicker", nil)
}
