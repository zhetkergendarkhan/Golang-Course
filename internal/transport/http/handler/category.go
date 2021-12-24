package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"shop/internal/model"
	"shop/internal/service"
	"strconv"
)

type CategoryHandler struct {
	cs service.CategoryService
}

func NewCategoryHandler(cs service.CategoryService) *CategoryHandler {
	return &CategoryHandler{cs: cs}
}

func (c CategoryHandler) FindCategories(rw http.ResponseWriter, r *http.Request) {
	categories := c.cs.FindAll()
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(categories); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c CategoryHandler) CreateCategory(rw http.ResponseWriter, r *http.Request) {
	categoryRequest := &model.CategoryRequest{}
	decodeErr := json.NewDecoder(r.Body).Decode(categoryRequest)
	rw.Header().Set("Content-Type", "application/json")
	if decodeErr != nil {
		http.Error(rw, "Incorrect request body", http.StatusInternalServerError)
		return
	}
	category, err := c.cs.CreateCategory(categoryRequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(rw).Encode(category); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (c CategoryHandler) UpdateCategory(rw http.ResponseWriter, r *http.Request) {
	categoryRequest := &model.CategoryRequest{}
	id, parseErr := strconv.Atoi(chi.URLParam(r, "id"))
	if parseErr != nil {
		http.Error(rw, "Incorrect id", http.StatusBadRequest)
		return
	}
	decodeErr := json.NewDecoder(r.Body).Decode(categoryRequest)
	rw.Header().Set("Content-Type", "application/json")
	if decodeErr != nil {
		http.Error(rw, "Incorrect request body", http.StatusInternalServerError)
		return
	}
	category, err := c.cs.UpdateCategory(categoryRequest, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(rw).Encode(category); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (c CategoryHandler) DeleteCategory(rw http.ResponseWriter, r *http.Request) {
	id, parseErr := strconv.Atoi(chi.URLParam(r, "id"))
	if parseErr != nil {
		http.Error(rw, "Incorrect id", http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	response, err := c.cs.DeleteCategory(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
