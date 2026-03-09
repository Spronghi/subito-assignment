package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/handler"
	"github.com/simonecolaci/subito-assignment/internal/repository"
	"github.com/simonecolaci/subito-assignment/internal/service"
)

func newTestOrderMux(t *testing.T) *http.ServeMux {
	t.Helper()

	db, erro := repository.NewSQLiteDB(":memory:")
	if erro != nil {
		t.Fatalf("failed to create in-memory database: %v", erro)
	}

	orderRepo, err := repository.NewSQLiteOrderRepository(db)
	if err != nil {
		t.Fatalf("failed to create order repository: %v", err)
	}

	productRepo, err := repository.NewSQLiteProductRepository(db)
	if err != nil {
		t.Fatalf("failed to create product repository: %v", err)
	}

	err = productRepo.Populate()
	if err != nil {
		t.Fatalf("failed to populate product repository: %v", err)
	}

	err = orderRepo.Populate()
	if err != nil {
		t.Fatalf("failed to populate order repository: %v", err)
	}

	orderService := service.NewOrderService(orderRepo, productRepo)

	mux := http.NewServeMux()

	handler.NewOrderHandler(orderService).RegisterRoutes(mux)

	return mux
}

func TestOrderHandler_Create_StatusOk(t *testing.T) {
	mux := newTestOrderMux(t)

	body := `{"items":[{"product_id":1,"quantity":2}]}`

	req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestOrderHandler_Create_Empty(t *testing.T) {
	mux := newTestOrderMux(t)

	req := httptest.NewRequest("POST", "/orders", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestOrderHandler_Create_InvalidProductID(t *testing.T) {
	mux := newTestOrderMux(t)

	body := `{"items":[{"product_id":-1,"quantity":2}]}`

	req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestOrderHandler_GetByID_Ok(t *testing.T) {
	mux := newTestOrderMux(t)

	req := httptest.NewRequest("GET", "/orders/1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var order entity.Order
	if err := json.NewDecoder(rec.Body).Decode(&order); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if order.ID != 1 {
		t.Errorf("expected order ID %d, got %d", 1, order.ID)
	}
}

func TestOrderHandler_GetByID_NotFound(t *testing.T) {
	mux := newTestOrderMux(t)

	req := httptest.NewRequest("GET", "/orders/0", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestOrderHandler_List(t *testing.T) {
	mux := newTestOrderMux(t)

	req := httptest.NewRequest("GET", "/orders", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var orders []entity.Order
	if err := json.NewDecoder(rec.Body).Decode(&orders); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("expected %d orders, got %d", 2, len(orders))
	}
}
