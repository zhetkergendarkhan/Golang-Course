package repository

import (
	"database/sql"
	"log"
	"shop/internal/model"
)

type ProductRepository interface {
	FindAllProductByCategory(categoryID int) []*model.Product
	FindProductByID(productID int) (*model.Product, error)
	CreateProduct(request *model.ProductRequest)
}

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepositoryImpl(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

func (p ProductRepositoryImpl) FindProductByID(productID int) (*model.Product, error) {
	return nil, nil
}

func (p ProductRepositoryImpl) FindAllProductByCategory(categoryID int) []*model.Product {
	products := make([]*model.Product, 0)
	rows, queryErr := p.db.Query("Select * from products where category_id=$1", categoryID)
	if queryErr != nil {
		return nil
	}
	for rows.Next() {
		product := &model.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil
		}
		products = append(products, product)
	}
	return products
}

func (p ProductRepositoryImpl) CreateProduct(request *model.ProductRequest) {
	_, execErr := p.db.Exec("INSERT INTO products (name, price, category_id) values ($1, $2, $3)", request.Name, request.Price, request.CategoryId)
	if execErr != nil {
		log.Println(execErr)
		return
	}
	log.Println("Successfully created!")

}
