package http

import (
	"shop/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

func configureRouter(c *chi.Mux, manager *handler.Manager) {
	c.Route("/api/v1", func(r chi.Router) {
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", manager.FindCategories)
			r.Post("/create", manager.CreateCategory)
			r.Put("/update/{id}", manager.UpdateCategory)
			r.Delete("/delete/{id}", manager.DeleteCategory)
			r.Get("/{id}/products", manager.FindAllProductByCategory)
		})
		r.Post("/product/create", manager.CreateProduct)
	})
}
