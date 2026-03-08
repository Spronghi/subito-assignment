package entity_test

import (
	"errors"
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
)

func TestProduct_VATAmount(t *testing.T) {
	p := &entity.Product{Price: 10000, VATRate: 0.22}

	vat, err := p.VATAmount()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if vat != 2200 {
		t.Errorf("expected VAT 2200, got %d", vat)
	}
}

func TestProduct_VATAmount_NegativePrice(t *testing.T) {
	p := &entity.Product{Price: -1, VATRate: 0.22}

	_, err := p.VATAmount()
	if !errors.Is(err, entity.DataInconsistency) {
		t.Errorf("expected DataInconsistency, got %v", err)
	}
}

func TestProduct_TotalPrice(t *testing.T) {
	p := &entity.Product{Price: 10000, VATRate: 0.22}

	total, err := p.TotalPrice()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 12200 {
		t.Errorf("expected total 12200, got %d", total)
	}
}

func TestProduct_TotalPrice_ZeroVAT(t *testing.T) {
	p := &entity.Product{Price: 10000, VATRate: 0}

	total, err := p.TotalPrice()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 10000 {
		t.Errorf("expected total 10000, got %d", total)
	}
}

func TestProduct_TotalPrice_NegativePrice(t *testing.T) {
	p := &entity.Product{Price: -1, VATRate: 0.22}

	_, err := p.TotalPrice()
	if !errors.Is(err, entity.DataInconsistency) {
		t.Errorf("expected DataInconsistency, got %v", err)
	}
}
