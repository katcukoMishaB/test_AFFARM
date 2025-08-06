package services

import (
	"context"
	"time"
	"tracker-core/internal/models"

	"github.com/uptrace/bun"
)

type PriceService struct {
	db *bun.DB
}

func NewPriceService(db *bun.DB) *PriceService {
	return &PriceService{db: db}
}

func (s *PriceService) SavePrice(ctx context.Context, currencyID int, price float64, timestamp time.Time) error {
	priceModel := &models.Price{
		CurrencyID: currencyID,
		Price:      price,
		Timestamp:  timestamp,
	}

	_, err := s.db.NewInsert().
		Model(priceModel).
		Exec(ctx)

	return err
}

func (s *PriceService) GetPriceAtTime(ctx context.Context, symbol string, timestamp time.Time) (*models.Price, error) {
	price := new(models.Price)

	err := s.db.NewSelect().
		Model(price).
		Join("JOIN currencies c ON c.id = price.currency_id").
		Where("c.symbol = ?", symbol).
		Order("ABS(EXTRACT(EPOCH FROM (price.timestamp - ?)))", timestamp.Format(time.RFC3339)).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return price, nil
}

func (s *PriceService) GetLatestPrice(ctx context.Context, symbol string) (*models.Price, error) {
	price := new(models.Price)

	err := s.db.NewSelect().
		Model(price).
		Join("JOIN currencies c ON c.id = price.currency_id").
		Where("c.symbol = ?", symbol).
		Order("price.timestamp DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return price, nil
}

func (s *PriceService) GetPriceHistory(ctx context.Context, symbol string, limit int) ([]models.Price, error) {
	var prices []models.Price

	err := s.db.NewSelect().
		Model(&prices).
		Join("JOIN currencies c ON c.id = price.currency_id").
		Where("c.symbol = ?", symbol).
		Order("price.timestamp DESC").
		Limit(limit).
		Scan(ctx)

	return prices, err
}

func (s *PriceService) GetPricesByCurrencyID(ctx context.Context, currencyID int) ([]models.Price, error) {
	var prices []models.Price

	err := s.db.NewSelect().
		Model(&prices).
		Where("currency_id = ?", currencyID).
		Order("timestamp DESC").
		Scan(ctx)

	return prices, err
}
