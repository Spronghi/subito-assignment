package entity_test

import (
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product *entity.Product
		wantErr bool
	}{
		{
			name: "valid product",
			product: &entity.Product{
				Name:        "Test Product",
				Description: "This is a test product",
				Price:       10000,
				VATRate:     0.22,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			product: &entity.Product{
				Description: "This is a test product",
				Price:       10000,
				VATRate:     0.22,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			product: &entity.Product{
				Name:        "Test Product",
				Description: "This is a test product",
				Price:       -1,
				VATRate:     0.22,
			},
			wantErr: true,
		},
		{
			name: "negative VAT rate",
			product: &entity.Product{
				Name:        "Test Product",
				Description: "This is a test product",
				Price:       10000,
				VATRate:     -0.01,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestProduct_VATAmount(t *testing.T) {
	p := &entity.Product{Price: 10000, VATRate: 0.22}

	if vat := p.VATAmount(); vat != 2200 {
		t.Errorf("expected VAT 2200, got %d", vat)
	}
}

func TestProduct_VATAmount_NegativePrice(t *testing.T) {
	p := &entity.Product{Price: -1, VATRate: 0.22}

	if vat := p.VATAmount(); vat != 0 {
		t.Errorf("expected VAT 0 for negative price, got %d", vat)
	}
}

func TestProduct_TotalPrice(t *testing.T) {
	p := &entity.Product{Price: 10000, VATRate: 0.22}

	if total := p.TotalPrice(); total != 12200 {
		t.Errorf("expected total 12200, got %d", total)
	}
}

func TestProduct_TotalPrice_ZeroVAT(t *testing.T) {
	p := &entity.Product{Price: 10000, VATRate: 0}

	if total := p.TotalPrice(); total != 10000 {
		t.Errorf("expected total 10000, got %d", total)
	}
}

func TestProduct_TotalPrice_NegativePrice(t *testing.T) {
	p := &entity.Product{Price: -1, VATRate: 0.22}

	if total := p.TotalPrice(); total != -1 {
		t.Errorf("expected total -1 for negative price, got %d", total)
	}
}
