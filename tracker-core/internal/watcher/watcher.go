package watcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"tracker-core/internal/models"
	"tracker-core/internal/services"

	"github.com/uptrace/bun"
)

type Watcher struct {
	db              *bun.DB
	currencyService *services.CurrencyService
	priceService    *services.PriceService
	stopChan        chan bool
	wg              sync.WaitGroup
	updateInterval  time.Duration
}

type KuCoinResponse struct {
	Code string `json:"code"`
	Data struct {
		Price string `json:"price"`
	} `json:"data"`
}

func NewWatcher(db *bun.DB, currencyService *services.CurrencyService, priceService *services.PriceService) *Watcher {
	return &Watcher{
		db:              db,
		currencyService: currencyService,
		priceService:    priceService,
		stopChan:        make(chan bool),
		updateInterval:  10 * time.Second,
	}
}

func (w *Watcher) Start() {
	w.wg.Add(1)
	go w.watchLoop()
	log.Println("Crypto watcher started")
}

func (w *Watcher) Stop() {
	close(w.stopChan)
	w.wg.Wait()
	log.Println("Crypto watcher stopped")
}

func (w *Watcher) watchLoop() {
	defer w.wg.Done()

	ticker := time.NewTicker(w.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.updatePrices()
		}
	}
}

func (w *Watcher) updatePrices() {
	ctx := context.Background()
	currencies, err := w.currencyService.GetActiveCurrencies(ctx)
	if err != nil {
		log.Printf("Error getting active currencies: %v", err)
		return
	}

	if len(currencies) == 0 {
		return
	}

	results := make(chan error, len(currencies))

	for _, currency := range currencies {
		go func(c models.Currency) {
			results <- w.updateCurrencyPrice(c)
		}(currency)
	}

	for i := 0; i < len(currencies); i++ {
		if err := <-results; err != nil {
			log.Printf("Price update error: %v", err)
		}
	}
}

func (w *Watcher) updateCurrencyPrice(currency models.Currency) error {
	price, err := w.fetchPriceFromKuCoin(currency.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get price for %s: %v", currency.Symbol, err)
	}

	ctx := context.Background()
	err = w.priceService.SavePrice(ctx, currency.ID, price, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save price for %s: %v", currency.Symbol, err)
	}

	log.Printf("Updated price %s: %.8f", currency.Symbol, price)
	return nil
}

func (w *Watcher) fetchPriceFromKuCoin(symbol string) (float64, error) {
	tradingPair := fmt.Sprintf("%s-USDT", strings.ToUpper(symbol))

	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=%s", tradingPair)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var kucoinResp KuCoinResponse
	if err := json.Unmarshal(body, &kucoinResp); err != nil {
		return 0, err
	}

	if kucoinResp.Code != "200000" {
		return 0, fmt.Errorf("API returned error: %s", kucoinResp.Code)
	}

	price, err := strconv.ParseFloat(kucoinResp.Data.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("price parsing error: %v", err)
	}

	return price, nil
}

func (w *Watcher) GetUpdateInterval() time.Duration {
	return w.updateInterval
}

func (w *Watcher) SetUpdateInterval(interval time.Duration) {
	w.updateInterval = interval
}