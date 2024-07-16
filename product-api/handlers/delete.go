package handlers

import (
	"microservices/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes product by id
// responses:
// 	201: noContentResponse

// DeleteProduct deletes the product from data store
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	// this will always convert because of the router
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle DELETE Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not found", http.StatusInternalServerError)
		return
	}
}
