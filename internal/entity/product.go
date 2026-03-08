package entity

import (
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	VATRate     float64   `json:"vat_rate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) VATAmount() int64 {
	// TODO: not sure how to handle negative prices, for now we'll just return 0
	if p.Price < 0 {
		return 0
	}

	return int64(float64(p.Price) * p.VATRate)
}

func (p *Product) TotalPrice() int64 {
	return p.Price + p.VATAmount()
}
