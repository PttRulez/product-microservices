package handlers

import (
	"net/http"

	"github.com/pttrulez/product-microservices/product_api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes product by id
// responses:
// 	201: noContentResponse

// DeleteProduct deletes the product from data store
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Debug("Handle DELETE Product", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("product not found", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("failed to delete product with id", id, "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
}
