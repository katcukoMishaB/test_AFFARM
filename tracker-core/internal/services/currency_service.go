package services

import (
	"context"
	"time"
	"tracker-core/internal/models"

	"github.com/uptrace/bun"
)

type CurrencyService struct {
	db *bun.DB
}

func NewCurrencyService(db *bun.DB) *CurrencyService {
	return &CurrencyService{db: db}
}

func (s *CurrencyService) AddCurrency(ctx context.Context, symbol, name string) error {
	currency := &models.Currency{
		Symbol:   symbol,
		Name:     name,
		IsActive: true,
	}

	_, err := s.db.NewInsert().
		Model(currency).
		On("CONFLICT (symbol) DO UPDATE").
		Set("is_active = EXCLUDED.is_active").
		Set("updated_at = now()").
		Exec(ctx)

	return err
}

func (s *CurrencyService) RemoveCurrency(ctx context.Context, symbol string) error {
	_, err := s.db.NewUpdate().
		Model((*models.Currency)(nil)).
		Set("is_active = ?", false).
		Set("updated_at = ?", time.Now()).
		Where("symbol = ?", symbol).
		Exec(ctx)

	return err
}

func (s *CurrencyService) GetActiveCurrencies(ctx context.Context) ([]models.Currency, error) {
	var currencies []models.Currency

	err := s.db.NewSelect().
		Model(&currencies).
		Where("is_active = ?", true).
		Order("symbol ASC").
		Scan(ctx)

	return currencies, err
}

func (s *CurrencyService) GetCurrencyBySymbol(ctx context.Context, symbol string) (*models.Currency, error) {
	currency := new(models.Currency)

	err := s.db.NewSelect().
		Model(currency).
		Where("symbol = ?", symbol).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return currency, nil
}

func (s *CurrencyService) GetCurrencyByID(ctx context.Context, id int) (*models.Currency, error) {
	currency := new(models.Currency)

	err := s.db.NewSelect().
		Model(currency).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return currency, nil
}
