package service_test

import (
	"errors"
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/repository"
	"github.com/simonecolaci/subito-assignment/internal/service"
)

func newTestProductService(t *testing.T) service.ProductService {
	t.Helper()

	db, err := repository.NewSQLiteDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	r, err := repository.NewSQLiteProductRepository(db)
	if err != nil {
		t.Fatalf("failed to create product repository: %v", err)
	}

	return service.NewProductService(r)
}

func TestProductService_Create(t *testing.T) {
	s := newTestProductService(t)

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
	s := newTestProductService(t)

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
	s := newTestProductService(t)

	_, err := s.Create(&entity.Product{Name: "Product 1", Description: "Description 1", Price: 10000, VATRate: 0.22})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

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
	s := newTestProductService(t)

	_, err := s.GetByID(0)
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected error %v, got %v", entity.ErrNotFound, err)
	}
}

func TestProductService_List(t *testing.T) {
	s := newTestProductService(t)

	productsToCreate := []entity.Product{
		{Name: "Product 1", Description: "Description 1", Price: 10000, VATRate: 0.22},
		{Name: "Product 2", Description: "Description 2", Price: 15000, VATRate: 0.22},
		{Name: "Product 3", Description: "Description 3", Price: 20000, VATRate: 0.22},
	}

	for i := range productsToCreate {
		if _, err := s.Create(&productsToCreate[i]); err != nil {
			t.Fatalf("failed to create product: %v", err)
		}
	}

	products, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}
}

func TestProductService_List_Empty(t *testing.T) {
	s := newTestProductService(t)

	products, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestProductService_Update(t *testing.T) {
	s := newTestProductService(t)

	created, err := s.Create(&entity.Product{Name: "Original", Description: "Desc", Price: 10000, VATRate: 0.22})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	updated, err := s.Update(created.ID, &entity.Product{Name: "Updated", Description: "New Desc", Price: 20000, VATRate: 0.10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Name != "Updated" {
		t.Errorf("expected name 'Updated', got '%s'", updated.Name)
	}
	if updated.Price != 20000 {
		t.Errorf("expected price 20000, got %d", updated.Price)
	}
}

func TestProductService_Update_NotFound(t *testing.T) {
	s := newTestProductService(t)

	_, err := s.Update(999, &entity.Product{Name: "Updated", Price: 10000, VATRate: 0.22})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected error %v, got %v", entity.ErrNotFound, err)
	}
}

func TestProductService_Update_InvalidInput(t *testing.T) {
	s := newTestProductService(t)

	created, err := s.Create(&entity.Product{Name: "Original", Description: "Desc", Price: 10000, VATRate: 0.22})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	_, err = s.Update(created.ID, &entity.Product{Name: "", Price: 10000, VATRate: 0.22})
	if !errors.Is(err, entity.ErrInvalidInput) {
		t.Fatalf("expected error %v, got %v", entity.ErrInvalidInput, err)
	}
}

func TestProductService_Delete(t *testing.T) {
	s := newTestProductService(t)

	created, err := s.Create(&entity.Product{Name: "ToDelete", Description: "Desc", Price: 10000, VATRate: 0.22})
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.GetByID(created.ID)
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected product to be deleted, got %v", err)
	}
}

func TestProductService_Delete_NotFound(t *testing.T) {
	s := newTestProductService(t)

	if err := s.Delete(999); !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected error %v, got %v", entity.ErrNotFound, err)
	}
}
