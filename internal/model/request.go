package model

type CategoryRequest struct {
	Name string `json:"name"`
}

type ProductRequest struct {
	Name       string  `json:"name"`
	Price      float32 `json:"price"`
	CategoryId int     `json:"categoryId"`
}
