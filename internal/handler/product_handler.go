package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products", h.list)
	mux.HandleFunc("GET /products/{id}", h.getByID)
	mux.HandleFunc("POST /products", h.create)
	mux.HandleFunc("PUT /products/{id}", h.update)
	mux.HandleFunc("DELETE /products/{id}", h.delete)
}

func (h *ProductHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if items == nil {
		items = []*entity.Product{}
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.service.GetByID(id)
	if errors.Is(err, entity.ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       int64   `json:"price"`
		VATRate     float64 `json:"vat_rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.service.Update(id, &entity.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		VATRate:     input.VATRate,
	})
	if errors.Is(err, entity.ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, entity.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.Delete(id); errors.Is(err, entity.ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       int64   `json:"price"`
		VATRate     float64 `json:"vat_rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.service.Create(&entity.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		VATRate:     input.VATRate,
	})
	if errors.Is(err, entity.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}
