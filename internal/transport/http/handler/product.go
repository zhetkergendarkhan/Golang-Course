package handler

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"net/http"
	"shop/internal/model"
	"shop/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ps service.ProductService
}

func NewProductHandler(ps service.ProductService) *ProductHandler {
	return &ProductHandler{ps: ps}
}

func (p ProductHandler) FindAllProductByCategory(rw http.ResponseWriter, r *http.Request) {
	id, parseErr := strconv.Atoi(chi.URLParam(r, "id"))
	if parseErr != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	products := p.ps.FindAllProductByCategory(id)
	if err := json.NewEncoder(rw).Encode(products); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p ProductHandler) CreateProduct(rw http.ResponseWriter, r *http.Request) {
	var productRequest []model.ProductRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&productRequest)
	rw.Header().Set("Content-Type", "application/json")
	if decodeErr != nil {
		http.Error(rw, decodeErr.Error(), http.StatusInternalServerError)
		return
	}
	response, err := p.ps.SendProduct(productRequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (p ProductHandler) ProductCreate(messages <-chan *message.Message) {
	productRequest := model.ProductRequest{}
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		unmarshalErr := json.Unmarshal(msg.Payload, &productRequest)
		if unmarshalErr != nil {
			log.Printf("received message json error")
			continue
		}
		p.ps.CreateProduct(&productRequest)
		msg.Ack()
	}
}
