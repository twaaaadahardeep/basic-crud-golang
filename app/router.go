package app

import (
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, service Service) {
	api := app.Group("/api")

	api.Get("/products", service.getProducts)
	api.Get("/product/:id", service.getProduct)
	api.Post("/product", service.createProduct)
	api.Put("/product/:id", service.updateProduct)
	api.Delete("/product/:id", service.deleteProduct)
}
