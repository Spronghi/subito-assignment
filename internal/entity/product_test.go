package entity_test

import (
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
)

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
