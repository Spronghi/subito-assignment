package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/service"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service}
}

func (h *OrderHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /orders", h.list)
	mux.HandleFunc("GET /orders/{id}", h.getByID)
	mux.HandleFunc("POST /orders", h.create)
	mux.HandleFunc("PUT /orders/{id}", h.update)
	mux.HandleFunc("DELETE /orders/{id}", h.delete)
}

func (h *OrderHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if items == nil {
		items = []*entity.Order{}
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *OrderHandler) getByID(w http.ResponseWriter, r *http.Request) {
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

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var input entity.NewOrder

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	order, err := h.service.Create(&input)
	if errors.Is(err, entity.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if errors.Is(err, entity.ErrNotFound) {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var input entity.NewOrder

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	order, err := h.service.Update(id, &input)
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

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) delete(w http.ResponseWriter, r *http.Request) {
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
