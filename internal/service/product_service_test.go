package service_test

import (
	"errors"
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/service"
)

func TestProductService_Create(t *testing.T) {
	s := service.NewProductService()

	p := &entity.Product{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       10000,
		VATRate:     0.22,
	}

	created, err := s.Create(p)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID == 0 {
		t.Errorf("expected product ID to be set, got %d", created.ID)
	}

	if created.Name != p.Name {
		t.Errorf("expected name '%s', got '%s'", p.Name, created.Name)
	}
}

func TestProductService_Create_InvalidInput(t *testing.T) {
	s := service.NewProductService()

	p := &entity.Product{
		Description: "This is a test product",
		Price:       10000,
		VATRate:     0.22,
	}

	_, err := s.Create(p)
	if !errors.Is(err, entity.ErrInvalidInput) {
		t.Fatalf("expected error %v, got %v", entity.ErrInvalidInput, err)
	}
}

func TestProductService_GetByID(t *testing.T) {
	s := service.NewProductService()

	p, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p.ID != 1 {
		t.Errorf("expected product ID 1, got %d", p.ID)
	}

	if p.Name != "Product 1" {
		t.Errorf("expected name 'Product 1', got '%s'", p.Name)
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	s := service.NewProductService()

	_, err := s.GetByID(0)
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected error %v, got %v", entity.ErrNotFound, err)
	}
}

func TestProductService_List(t *testing.T) {
	s := service.NewProductService()

	products, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}
}
