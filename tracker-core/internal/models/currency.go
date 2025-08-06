package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Currency struct {
	bun.BaseModel `bun:"table:currencies"`

	ID        int       `bun:"id,pk,autoincrement" json:"id"`
	Symbol    string    `bun:"symbol,unique,notnull" json:"symbol"`
	Name      string    `bun:"name,notnull" json:"name"`
	IsActive  bool      `bun:"is_active,default:true" json:"is_active"`
	CreatedAt time.Time `bun:"created_at,default:now()" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,default:now()" json:"updated_at"`

	Prices []Price `bun:"rel:has-many,join:id=currency_id" json:"prices,omitempty"`
}

func (Currency) TableName() string {
	return "currencies"
}
