package entity_test

import (
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
)

func TestOrder_Validate(t *testing.T) {
	tests := []struct {
		name    string
		order   *entity.Order
		wantErr bool
	}{
		{
			name: "valid order",
			order: &entity.Order{
				Items: []entity.OrderItem{
					{ProductID: 1, Quantity: 2},
					{ProductID: 2, Quantity: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "empty items",
			order: &entity.Order{
				Items: []entity.OrderItem{},
			},
			wantErr: true,
		},
		{
			name: "negative quantity",
			order: &entity.Order{
				Items: []entity.OrderItem{
					{ProductID: 1, Quantity: -1},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestOrder_NewOrder_Validate(t *testing.T) {
	tests := []struct {
		name    string
		order   *entity.NewOrder
		wantErr bool
	}{
		{
			name: "valid new order",
			order: &entity.NewOrder{
				Items: []entity.NewOrderItem{
					{ProductID: 1, Quantity: 2},
					{ProductID: 2, Quantity: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "empty items",
			order: &entity.NewOrder{
				Items: []entity.NewOrderItem{},
			},
			wantErr: true,
		},
		{
			name: "negative quantity",
			order: &entity.NewOrder{
				Items: []entity.NewOrderItem{
					{ProductID: 1, Quantity: -1},
				},
			},
			wantErr: true,
		},
		{
			name: "zero quantity",
			order: &entity.NewOrder{
				Items: []entity.NewOrderItem{
					{ProductID: 1, Quantity: 0},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestOrder_CalculateTotalPrice(t *testing.T) {
	order := &entity.Order{
		Items: []entity.OrderItem{
			{UnitPrice: 10000, Quantity: 2, Price: 20000},
			{UnitPrice: 5000, Quantity: 1, Price: 5000},
		},
	}

	expectedTotal := int64(25000)
	if order.CalculateTotalPrice() != expectedTotal {
		t.Errorf("expected total price %d, got %d", expectedTotal, order.CalculateTotalPrice())
	}
}

func TestOrder_CalculateTotalVAT(t *testing.T) {
	order := &entity.Order{
		Items: []entity.OrderItem{
			{VAT: 2200},
			{VAT: 1100},
		},
	}

	expectedTotalVAT := int64(3300)
	if order.CalculateTotalVAT() != expectedTotalVAT {
		t.Errorf("expected total VAT %d, got %d", expectedTotalVAT, order.CalculateTotalVAT())
	}
}
