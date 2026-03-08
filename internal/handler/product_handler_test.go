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

func newTestProductMux(t *testing.T) *http.ServeMux {
	t.Helper()

	db, err := repository.NewSQLiteDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	productRepo, err := repository.NewSQLiteProductRepository(db)
	if err != nil {
		t.Fatalf("failed to create product repository: %v", err)
	}

	if err := productRepo.Populate(); err != nil {
		t.Fatalf("failed to populate product repository: %v", err)
	}

	productService := service.NewProductService(productRepo)

	mux := http.NewServeMux()

	handler.NewProductHandler(productService).RegisterRoutes(mux)

	return mux
}

func TestProductHandler_ListOkStatus(t *testing.T) {
	mux := newTestProductMux(t)

	req := httptest.NewRequest("GET", "/products", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestProductHandler_List(t *testing.T) {
	mux := newTestProductMux(t)

	req := httptest.NewRequest("GET", "/products", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var products []entity.Product
	if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(products) != 3 {
		t.Fatalf("expected %d products, got %d", 3, len(products))
	}
}

func TestProductHandler_GetByID(t *testing.T) {
	mux := newTestProductMux(t)

	req := httptest.NewRequest("GET", "/products/1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var product entity.Product
	if err := json.NewDecoder(rec.Body).Decode(&product); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if product.ID != 1 {
		t.Errorf("expected product ID %d, got %d", 1, product.ID)
	}
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	mux := newTestProductMux(t)

	req := httptest.NewRequest("GET", "/products/0", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestProductHandler_Create(t *testing.T) {
	mux := newTestProductMux(t)

	body := `{"name":"New Product","description":"A new product","price":15000,"vat_rate":0.22}`
	req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var product entity.Product
	if err := json.NewDecoder(rec.Body).Decode(&product); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if product.Name != "New Product" {
		t.Errorf("expected product name '%s', got '%s'", "New Product", product.Name)
	}
}

func TestProductHandler_Create_BadRequest(t *testing.T) {
	mux := newTestProductMux(t)

	body := `{"name":"","description":"A new product","price":15000,"vat_rate":0.22}`
	req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
