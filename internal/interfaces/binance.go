package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/constants"
	"github.com/ntdat104/go-finance-dataset/internal/application/response"
	"github.com/ntdat104/go-finance-dataset/internal/application/service"
)

type BinanceHandler interface {
	Ping(ctx *gin.Context)
	ServerTime(ctx *gin.Context)
	ExchangeInfo(ctx *gin.Context)
	TickerPrice(ctx *gin.Context)
	AllPrices(ctx *gin.Context)
	BookTicker(ctx *gin.Context)
	Depth(ctx *gin.Context)
	RecentTrades(ctx *gin.Context)
	Klines(ctx *gin.Context)
	HistoricalTrades(ctx *gin.Context)
	AggregateTrades(ctx *gin.Context)
	AvgPrice(ctx *gin.Context)
	Ticker24Hr(ctx *gin.Context)
	AllBookTickers(ctx *gin.Context)
}

type binanceHandler struct {
	router     *gin.Engine
	binanceSvc service.BinanceSvc
}

func NewBinanceHandler(router *gin.Engine, binanceSvc service.BinanceSvc) BinanceHandler {
	h := &binanceHandler{
		router:     router,
		binanceSvc: binanceSvc,
	}
	h.initRoutes()
	return h
}

func (h *binanceHandler) initRoutes() {
	h.router.GET(constants.ApiBinancePing, h.Ping)
	h.router.GET(constants.ApiBinanceServerTime, h.ServerTime)
	h.router.GET(constants.ApiBinanceExchangeInfo, h.ExchangeInfo)
	h.router.GET(constants.ApiBinanceTickerPrice, h.TickerPrice)
	h.router.GET(constants.ApiBinanceAllPrices, h.AllPrices)
	h.router.GET(constants.ApiBinanceBookTicker, h.BookTicker)
	h.router.GET(constants.ApiBinanceDepth, h.Depth)
	h.router.GET(constants.ApiBinanceRecentTrades, h.RecentTrades)
	h.router.GET(constants.ApiBinanceKlines, h.Klines)
	h.router.GET(constants.ApiBinanceHistoricalTrades, h.HistoricalTrades)
	h.router.GET(constants.ApiBinanceAggregateTrades, h.AggregateTrades)
	h.router.GET(constants.ApiBinanceAvgPrice, h.AvgPrice)
	h.router.GET(constants.ApiBinanceTicker24Hr, h.Ticker24Hr)
	h.router.GET(constants.ApiBinanceAllBookTickers, h.AllBookTickers)
}

// Ping handles the /api/v3/ping endpoint.
func (c *binanceHandler) Ping(ctx *gin.Context) {
	resp, err := c.binanceSvc.GetPing()
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// ServerTime handles the /api/v3/time endpoint.
func (c *binanceHandler) ServerTime(ctx *gin.Context) {
	resp, err := c.binanceSvc.GetServerTime()
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// ExchangeInfo handles the /api/v3/exchangeInfo endpoint.
func (c *binanceHandler) ExchangeInfo(ctx *gin.Context) {
	resp, err := c.binanceSvc.GetExchangeInfo()
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// TickerPrice handles the /api/v3/ticker/price endpoint for a single symbol.
func (c *binanceHandler) TickerPrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceSvc.GetTickerPrice(symbol)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// AllPrices handles the /api/v3/ticker/price endpoint for all symbols.
func (c *binanceHandler) AllPrices(ctx *gin.Context) {
	resp, err := c.binanceSvc.GetAllTickerPrices()
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// BookTicker handles the /api/v3/ticker/bookTicker endpoint for a single symbol.
func (c *binanceHandler) BookTicker(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceSvc.GetBookTicker(symbol)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// Depth handles the /api/v3/depth endpoint.
func (c *binanceHandler) Depth(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceSvc.GetDepth(symbol, limit)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// RecentTrades handles the /api/v3/trades endpoint.
func (c *binanceHandler) RecentTrades(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceSvc.GetRecentTrades(symbol, limit)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// Klines handles the /api/v3/klines endpoint.
func (c *binanceHandler) Klines(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	interval := ctx.Query("interval")
	if symbol == "" || interval == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol and interval query parameters are required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceSvc.GetKlines(symbol, interval, limit)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// HistoricalTrades handles the /api/v3/historicalTrades endpoint.
func (c *binanceHandler) HistoricalTrades(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "500") // Default limit for historical trades
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	var fromId *int64
	fromIdStr := ctx.Query("fromId")
	if fromIdStr != "" {
		id, err := strconv.ParseInt(fromIdStr, 10, 64)
		if err != nil {
			response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid fromId parameter"})
			return
		}
		fromId = &id
	}

	resp, err := c.binanceSvc.GetHistoricalTrades(symbol, limit, fromId)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// AggregateTrades handles the /api/v3/aggTrades endpoint.
func (c *binanceHandler) AggregateTrades(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}

	var fromId, startTime, endTime *int64
	limit := 500 // Default limit

	if s := ctx.Query("fromId"); s != "" {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid fromId parameter"})
			return
		}
		fromId = &id
	}
	if s := ctx.Query("startTime"); s != "" {
		t, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid startTime parameter"})
			return
		}
		startTime = &t
	}
	if s := ctx.Query("endTime"); s != "" {
		t, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid endTime parameter"})
			return
		}
		endTime = &t
	}
	if s := ctx.Query("limit"); s != "" {
		l, err := strconv.Atoi(s)
		if err != nil {
			response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
		limit = l
	}

	resp, err := c.binanceSvc.GetAggregateTrades(symbol, fromId, startTime, endTime, limit)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// AvgPrice handles the /api/v3/avgPrice endpoint.
func (c *binanceHandler) AvgPrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceSvc.GetAvgPrice(symbol)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// Ticker24Hr handles the /api/v3/ticker/24hr endpoint.
func (c *binanceHandler) Ticker24Hr(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		response.JSON(ctx, http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceSvc.GetTicker24Hr(symbol)
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}

// AllBookTickers handles the /api/v3/ticker/bookTicker endpoint for all symbols.
func (c *binanceHandler) AllBookTickers(ctx *gin.Context) {
	resp, err := c.binanceSvc.GetAllBookTickers()
	if err != nil {
		response.JSON(ctx, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, resp)
}
