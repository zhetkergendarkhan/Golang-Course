package service

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"shop/internal/model"
	"shop/internal/repository"
)

type ProductService interface {
	FindAllProductByCategory(categoryID int) []*model.Product
	FindProductByID(productID int) (*model.Product, error)
	CreateProduct(request *model.ProductRequest)
	SendProduct(request []model.ProductRequest) (*model.SuccessResponse, error)
}

type ProductServiceImpl struct {
	pr        repository.ProductRepository
	publisher message.Publisher
}

func NewProductServiceImpl(pr repository.ProductRepository, publisher message.Publisher) ProductService {
	return &ProductServiceImpl{pr: pr, publisher: publisher}
}

func (p ProductServiceImpl) FindProductByID(productID int) (*model.Product, error) {
	return p.pr.FindProductByID(productID)
}

func (p ProductServiceImpl) FindAllProductByCategory(categoryID int) []*model.Product {
	return p.pr.FindAllProductByCategory(categoryID)
}

func (p ProductServiceImpl) CreateProduct(request *model.ProductRequest) {
	p.pr.CreateProduct(request)
	return
}

func (p ProductServiceImpl) SendProduct(request []model.ProductRequest) (*model.SuccessResponse, error) {
	for _, value := range request {
		jsonData, _ := json.Marshal(value)
		msg := message.NewMessage(watermill.NewUUID(), jsonData)
		publisherErr := p.publisher.Publish("product-create", msg)
		if publisherErr != nil {
			return nil, publisherErr
		}
	}

	return &model.SuccessResponse{Message: "Successfully send to broker"}, nil
}
