package entity

import (
	"fmt"
	"time"
)

type NewOrder struct {
	Items []NewOrderItem `json:"items"`
}

type NewOrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type Order struct {
	ID         int64       `json:"id"`
	Items      []OrderItem `json:"items"`
	TotalPrice int64       `json:"total_price"`
	TotalVAT   int64       `json:"total_vat"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID          int64   `json:"id,omitempty"`
	OrderID     int64   `json:"order_id,omitempty"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   int64   `json:"unit_price"`
	VATRate     float64 `json:"vat_rate"`
	Price       int64   `json:"price"`
	VAT         int64   `json:"vat"`
}

func (o Order) Validate() error {
	if len(o.Items) == 0 {
		return fmt.Errorf("%w: order must have at least one item", ErrInvalidInput)
	}
	for _, item := range o.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("invalid order item: %w", err)
		}
	}
	return nil
}

func (oi OrderItem) Validate() error {
	if oi.ProductID == 0 {
		return fmt.Errorf("product ID must be set: %w", ErrInvalidInput)
	}
	if oi.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0: %w", ErrInvalidInput)
	}
	return nil
}

func (o Order) CalculateTotalPrice() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price
	}
	return total
}

func (o Order) CalculateTotalVAT() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.VAT
	}
	return total
}

// MergeItems merges order items with the same product ID by summing their quantities and recalculating price and VAT.
func (oi Order) MergeItems() []OrderItem {
	merged := make(map[int64]OrderItem)

	for _, item := range oi.Items {
		if existing, ok := merged[item.ProductID]; ok {
			existing.Quantity += item.Quantity
			merged[item.ProductID] = existing
		} else {
			merged[item.ProductID] = item
		}
	}

	var mergedItems []OrderItem
	for _, item := range merged {
		mergedItems = append(mergedItems, item)
	}
	return mergedItems
}
