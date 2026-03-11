package service_test

import (
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/repository"
	"github.com/simonecolaci/subito-assignment/internal/service"
)

func newTestOrderService(t *testing.T) service.OrderService {
	t.Helper()

	db, err := repository.NewSQLiteDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	orderRepo, err := repository.NewSQLiteOrderRepository(db)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	productRepo, err := repository.NewSQLiteProductRepository(db)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	err = productRepo.Populate()
	if err != nil {
		t.Fatalf("failed to populate repository: %v", err)
	}

	return service.NewOrderService(orderRepo, productRepo)
}

func TestOrderService_Create(t *testing.T) {
	s := newTestOrderService(t)

	newOrder := &entity.NewOrder{
		Items: []entity.NewOrderItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	created, err := s.Create(newOrder)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID == 0 {
		t.Errorf("expected order ID to be set, got %d", created.ID)
	}

	if created.TotalPrice != 210000 {
		t.Errorf("expected total price 160000, got %d", created.TotalPrice)
	}

	if created.TotalVAT != 46200 {
		t.Errorf("expected total VAT 35200, got %d", created.TotalVAT)
	}

	if len(created.Items) != len(newOrder.Items) {
		t.Errorf("expected %d items, got %d", len(newOrder.Items), len(created.Items))
	}

	for i, item := range created.Items {
		expected := newOrder.Items[i]
		if item.ProductID != expected.ProductID {
			t.Errorf("expected product ID %d, got %d", expected.ProductID, item.ProductID)
		}
		if item.Quantity != expected.Quantity {
			t.Errorf("expected quantity %d, got %d", expected.Quantity, item.Quantity)
		}
	}
}

func TestOrderService_Create_InvalidInput(t *testing.T) {
	s := newTestOrderService(t)

	newOrder := &entity.NewOrder{
		Items: []entity.NewOrderItem{
			{ProductID: 1, Quantity: 0},
		},
	}

	_, err := s.Create(newOrder)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderService_GetByID(t *testing.T) {
	s := newTestOrderService(t)

	newOrder := &entity.NewOrder{
		Items: []entity.NewOrderItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	created, err := s.Create(newOrder)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	fetched, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, fetched.ID)
	}

	if len(fetched.Items) != len(newOrder.Items) {
		t.Errorf("expected %d items, got %d", len(newOrder.Items), len(fetched.Items))
	}

	for i, item := range fetched.Items {
		expected := newOrder.Items[i]

		if item.ProductID != expected.ProductID {
			t.Errorf("expected product ID %d, got %d", expected.ProductID, item.ProductID)
		}
		if item.Quantity != expected.Quantity {
			t.Errorf("expected quantity %d, got %d", expected.Quantity, item.Quantity)
		}
	}
}

func TestOrderService_GetByID_NotFound(t *testing.T) {
	s := newTestOrderService(t)

	_, err := s.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderService_List(t *testing.T) {
	s := newTestOrderService(t)

	ordersToCreate := []*entity.NewOrder{
		{
			Items: []entity.NewOrderItem{
				{ProductID: 1, Quantity: 2},
				{ProductID: 2, Quantity: 1},
			},
		},
		{
			Items: []entity.NewOrderItem{
				{ProductID: 3, Quantity: 1},
			},
		},
	}

	for _, o := range ordersToCreate {
		if _, err := s.Create(o); err != nil {
			t.Fatalf("failed to create order: %v", err)
		}
	}

	fetchedOrders, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(fetchedOrders) != len(ordersToCreate) {
		t.Errorf("expected %d orders, got %d", len(ordersToCreate), len(fetchedOrders))
	}
}

func TestOrderService_List_Empty(t *testing.T) {
	s := newTestOrderService(t)

	fetchedOrders, err := s.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(fetchedOrders) != 0 {
		t.Errorf("expected 0 orders, got %d", len(fetchedOrders))
	}
}

func TestOrderService_Create_ProductNotFound(t *testing.T) {
	s := newTestOrderService(t)

	newOrder := &entity.NewOrder{
		Items: []entity.NewOrderItem{
			{ProductID: 999, Quantity: 1},
		},
	}

	_, err := s.Create(newOrder)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderService_Create_ProductInvalid(t *testing.T) {
	s := newTestOrderService(t)

	newOrder := &entity.NewOrder{
		Items: []entity.NewOrderItem{
			{ProductID: 1, Quantity: -1},
		},
	}

	_, err := s.Create(newOrder)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderService_Create_OrderInvalid(t *testing.T) {
	s := newTestOrderService(t)

	newOrder := &entity.NewOrder{
		Items: []entity.NewOrderItem{},
	}

	_, err := s.Create(newOrder)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderService_Create_OrderItemInvalid(t *testing.T) {
	s := newTestOrderService(t)

	newOrder := &entity.NewOrder{
		Items: []entity.NewOrderItem{
			{ProductID: 1, Quantity: 0},
		},
	}

	_, err := s.Create(newOrder)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderService_Update(t *testing.T) {
	s := newTestOrderService(t)

	created, err := s.Create(&entity.NewOrder{
		Items: []entity.NewOrderItem{{ProductID: 1, Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	updated, err := s.Update(created.ID, &entity.NewOrder{
		Items: []entity.NewOrderItem{{ProductID: 2, Quantity: 3}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(updated.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(updated.Items))
	}
	if updated.Items[0].ProductID != 2 {
		t.Errorf("expected product ID 2, got %d", updated.Items[0].ProductID)
	}
	if updated.Items[0].Quantity != 3 {
		t.Errorf("expected quantity 3, got %d", updated.Items[0].Quantity)
	}
}

func TestOrderService_Update_NotFound(t *testing.T) {
	s := newTestOrderService(t)

	_, err := s.Update(999, &entity.NewOrder{
		Items: []entity.NewOrderItem{{ProductID: 1, Quantity: 1}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderService_Delete(t *testing.T) {
	s := newTestOrderService(t)

	created, err := s.Create(&entity.NewOrder{
		Items: []entity.NewOrderItem{{ProductID: 1, Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err == nil {
		t.Fatal("expected order to be deleted, got nil error")
	}
}

func TestOrderService_Delete_NotFound(t *testing.T) {
	s := newTestOrderService(t)

	if err := s.Delete(999); err == nil {
		t.Fatal("expected error, got nil")
	}
}
