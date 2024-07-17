package handlers

import (
	"net/http"

	"github.com/pttrulez/product-microservices/product_api/data"
)

func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Debug("Updating record id", prod.ID)

	err := p.productDB.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Error("Product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
