package repository

import (
	"database/sql"
	"shop/internal/model"
)

type CategoryRepository interface {
	FindAll() []*model.Category
	CreateCategory(request *model.CategoryRequest) (*model.Category, error)
	UpdateCategory(request *model.CategoryRequest, id int) (*model.Category, error)
	DeleteCategory(id int) error
}

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepositoryImpl(db *sql.DB) CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (c CategoryRepositoryImpl) FindAll() []*model.Category {
	categories := make([]*model.Category, 0)
	rows, queryErr := c.db.Query("Select * from categories")
	if queryErr != nil {
		return nil
	}
	for rows.Next() {
		category := &model.Category{}
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil
		}
		categories = append(categories, category)
	}
	return categories
}

func (c CategoryRepositoryImpl) CreateCategory(request *model.CategoryRequest) (*model.Category, error) {
	category := &model.Category{}
	queryErr := c.db.QueryRow("INSERT INTO categories (name) values ($1) returning id, name", request.Name).Scan(&category.ID, &category.Name)
	if queryErr != nil {
		return nil, queryErr
	}
	return category, nil
}

func (c CategoryRepositoryImpl) UpdateCategory(request *model.CategoryRequest, id int) (*model.Category, error) {
	category := &model.Category{}
	queryErr := c.db.QueryRow("UPDATE categories set name=$1 where id=$2 returning id, name", request.Name, id).Scan(&category.ID, &category.Name)
	if queryErr != nil {
		return nil, queryErr
	}
	return category, nil
}

func (c CategoryRepositoryImpl) DeleteCategory(id int) error {
	tx, txErr := c.db.Begin()
	if txErr != nil {
		return txErr
	}
	_, execErr := tx.Exec("Delete from products where category_id=$1", id)
	if execErr != nil {
		tx.Rollback()
		return execErr
	}
	_, execErr = tx.Exec("Delete from categories where id=$1", id)
	if execErr != nil {
		tx.Rollback()
		return execErr
	}
	tx.Commit()
	return nil
}
