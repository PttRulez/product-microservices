package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pttrulez/product-microservices/product_api/data"
)

func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert Id", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
