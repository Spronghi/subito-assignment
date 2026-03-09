package repository_test

import (
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/repository"
)

func newTestOrderRepository(t *testing.T) repository.OrderRepository {
	t.Helper()

	db, err := repository.NewSQLiteDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	r, err := repository.NewSQLiteOrderRepository(db)
	if err != nil {
		t.Fatalf("failed to create order repository: %v", err)
	}

	return r
}

func TestSQLiteOrderRepository_Create(t *testing.T) {
	r := newTestOrderRepository(t)

	order := &entity.Order{
		Items: []entity.OrderItem{
			{ProductID: 1, Quantity: 2, UnitPrice: 10000, Price: 20000},
			{ProductID: 2, Quantity: 1, UnitPrice: 5000, Price: 5000},
		},
		TotalPrice: 25000,
	}

	created, err := r.Create(order)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID == 0 {
		t.Errorf("expected order ID to be set, got %d", created.ID)
	}

	if len(created.Items) != len(order.Items) {
		t.Errorf("expected %d items, got %d", len(order.Items), len(created.Items))
	}

	if created.TotalPrice != order.TotalPrice {
		t.Errorf("expected total price %d, got %d", order.TotalPrice, created.TotalPrice)
	}
}

func TestSQLiteOrderRepository_GetByID(t *testing.T) {
	r := newTestOrderRepository(t)

	order := &entity.Order{
		Items: []entity.OrderItem{
			{ProductID: 1, Quantity: 2, UnitPrice: 10000, Price: 20000},
		},
		TotalPrice: 20000,
	}

	created, err := r.Create(order)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	fetched, err := r.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, fetched.ID)
	}

	if len(fetched.Items) != len(order.Items) {
		t.Errorf("expected %d items, got %d", len(order.Items), len(fetched.Items))
	}

	if fetched.TotalPrice != order.TotalPrice {
		t.Errorf("expected total price %d, got %d", order.TotalPrice, fetched.TotalPrice)
	}
}

func TestSQLiteOrderRepository_GetByID_NotFound(t *testing.T) {
	r := newTestOrderRepository(t)

	_, err := r.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSQLiteOrderRepository_List(t *testing.T) {
	r := newTestOrderRepository(t)

	order1 := &entity.Order{
		Items: []entity.OrderItem{
			{ProductID: 1, Quantity: 2, UnitPrice: 10000, Price: 20000},
		},
		TotalPrice: 20000,
	}

	order2 := &entity.Order{
		Items: []entity.OrderItem{
			{ProductID: 2, Quantity: 1, UnitPrice: 5000, Price: 5000},
		},
		TotalPrice: 5000,
	}

	if _, err := r.Create(order1); err != nil {
		t.Fatalf("failed to create order1: %v", err)
	}

	if _, err := r.Create(order2); err != nil {
		t.Fatalf("failed to create order2: %v", err)
	}

	orders, err := r.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}
