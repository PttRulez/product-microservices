package handlers

import (
	"encoding/json"
	"microservices/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

// GetProducts return the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch list of products from datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	rw.Header().Add("Content-Type", "application/json")
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Println("[DEBUG] get record id", id)

	product, _, err := data.GetProductById(id)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching producr", err)

		http.Error(rw, "No product with id "+string(id), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(product)
	if err != nil {
		p.l.Println("[ERROR] fetching producr", err)

		http.Error(rw, "Failed to marshal product", http.StatusBadRequest)
		return
	}

	rw.Write(j)
}
