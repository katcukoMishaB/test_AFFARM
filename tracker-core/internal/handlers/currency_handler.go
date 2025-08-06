package handlers

import (
	"net/http"
	"time"
	"tracker-core/internal/services"

	"github.com/gin-gonic/gin"
)

type CurrencyHandler struct {
	currencyService *services.CurrencyService
	priceService    *services.PriceService
}

func NewCurrencyHandler(currencyService *services.CurrencyService, priceService *services.PriceService) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
		priceService:    priceService,
	}
}

type AddCurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type RemoveCurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required"`
}

type PriceRequest struct {
	Coin      string `json:"coin" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

func (h *CurrencyHandler) AddCurrency(c *gin.Context) {
	var req AddCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid data format",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	err := h.currencyService.AddCurrency(ctx, req.Symbol, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add currency",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Currency successfully added",
		"symbol":  req.Symbol,
		"name":    req.Name,
	})
}

func (h *CurrencyHandler) RemoveCurrency(c *gin.Context) {
	var req RemoveCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid data format",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	err := h.currencyService.RemoveCurrency(ctx, req.Symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove currency",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Currency successfully removed",
		"symbol":  req.Symbol,
	})
}

func (h *CurrencyHandler) GetPrice(c *gin.Context) {
	var req PriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid data format",
			"details": err.Error(),
		})
		return
	}

	timestamp := time.Unix(req.Timestamp, 0)
	ctx := c.Request.Context()

	price, err := h.priceService.GetPriceAtTime(ctx, req.Coin, timestamp)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Price not found",
			"coin":  req.Coin,
			"time":  timestamp,
		})
		return
	}

	c.JSON(http.StatusOK, price)
}

func (h *CurrencyHandler) GetActiveCurrencies(c *gin.Context) {
	ctx := c.Request.Context()
	currencies, err := h.currencyService.GetActiveCurrencies(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get currency list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, currencies)
}
