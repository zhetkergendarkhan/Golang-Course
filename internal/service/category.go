package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"shop/internal/model"
	"shop/internal/repository"
	"time"
)

type CategoryService interface {
	FindAll() []*model.Category
	CreateCategory(request *model.CategoryRequest) (*model.Category, error)
	UpdateCategory(request *model.CategoryRequest, id int) (*model.Category, error)
	DeleteCategory(id int) (*model.SuccessResponse, error)
}

const (
	categoriesKey = "categories"
)

type CategoryServiceImpl struct {
	cr  repository.CategoryRepository
	rdb *redis.Client
}

func NewCategoryServiceImpl(cr repository.CategoryRepository, rdb *redis.Client) CategoryService {
	return &CategoryServiceImpl{cr: cr, rdb: rdb}
}

func (c CategoryServiceImpl) FindAll() []*model.Category {
	categories := make([]*model.Category, 0)
	result, err := c.rdb.Get(context.Background(), categoriesKey).Result()
	if err != nil {
		if err == redis.Nil {
			jsonData, _ := json.Marshal(c.cr.FindAll())
			setErr := c.rdb.Set(context.Background(), categoriesKey, jsonData, 1*time.Minute).Err()
			if setErr != nil {
				log.Println(setErr)
			}
			return c.cr.FindAll()
		}
		return c.cr.FindAll()
	}
	unmarshalErr := json.Unmarshal([]byte(result), &categories)
	if unmarshalErr != nil {
		return c.cr.FindAll()
	}
	return categories
}

func (c CategoryServiceImpl) CreateCategory(request *model.CategoryRequest) (*model.Category, error) {
	return c.cr.CreateCategory(request)
}

func (c CategoryServiceImpl) UpdateCategory(request *model.CategoryRequest, id int) (*model.Category, error) {
	return c.cr.UpdateCategory(request, id)
}

func (c CategoryServiceImpl) DeleteCategory(id int) (*model.SuccessResponse, error) {
	return &model.SuccessResponse{Message: "Successfully deleted"}, c.cr.DeleteCategory(id)
}
