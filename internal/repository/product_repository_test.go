package repository_test

import (
	"errors"
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/repository"
)

func newTestProductRepository(t *testing.T) repository.ProductRepository {
	t.Helper()

	db, err := repository.NewSQLiteDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	r, err := repository.NewSQLiteProductRepository(db)
	if err != nil {
		t.Fatalf("failed to create product repository: %v", err)
	}

	return r
}

func TestSQLiteProductRepository_Create(t *testing.T) {
	r := newTestProductRepository(t)

	p := &entity.Product{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       10000,
		VATRate:     0.22,
	}

	created, err := r.Create(p)
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

func TestSQLiteProductRepository_List(t *testing.T) {
	r := newTestProductRepository(t)

	products := []*entity.Product{
		{Name: "Product 1", Description: "Description 1", Price: 10000, VATRate: 0.22},
		{Name: "Product 2", Description: "Description 2", Price: 20000, VATRate: 0.22},
	}

	for _, p := range products {
		if _, err := r.Create(p); err != nil {
			t.Fatalf("failed to create product: %v", err)
		}
	}

	listed, err := r.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(listed) != len(products) {
		t.Errorf("expected %d products, got %d", len(products), len(listed))
	}
}

func TestSQLiteProductRepository_GetByID(t *testing.T) {
	r := newTestProductRepository(t)

	p := &entity.Product{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       10000,
		VATRate:     0.22,
	}

	created, err := r.Create(p)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	fetched, err := r.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, fetched.ID)
	}
}

func TestSQLiteProductRepository_GetByID_NotFound(t *testing.T) {
	r := newTestProductRepository(t)

	_, err := r.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}
