package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Price struct {
	bun.BaseModel `bun:"table:prices"`

	ID         int       `bun:"id,pk,autoincrement" json:"id"`
	CurrencyID int       `bun:"currency_id,notnull" json:"currency_id"`
	Price      float64   `bun:"price,notnull" json:"price"`
	Timestamp  time.Time `bun:"timestamp,notnull" json:"timestamp"`
	CreatedAt  time.Time `bun:"created_at,default:now()" json:"created_at"`

	Currency *Currency `bun:"rel:belongs-to,join:currency_id=id" json:"currency,omitempty"`
}

func (Price) TableName() string {
	return "prices"
}
