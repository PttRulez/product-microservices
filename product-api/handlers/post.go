package handlers

import (
	"net/http"

	"github.com/pttrulez/product-microservices/product_api/data"
)

// swagger:route POST /products createProduct
// Create a new product
//
// responses:
// 200: productResponse
// 422: errorValidation
// 501: errorResponse

// Create handlers POST requests to add new products
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)
}
