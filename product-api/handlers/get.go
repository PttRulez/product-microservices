package handlers

import (
	"net/http"

	"github.com/pttrulez/product-microservices/product_api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

// GetProducts return the products from the data store
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Handle GET Products")
	rw.Header().Add("Content-Type", "application/json")
	cur := r.URL.Query().Get("currency")

	// fetch list of products from datastore
	lp, err := p.productDB.GetProducts(cur)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// serialize the list to JSON
	err = data.ToJSON(lp, rw)
	if err != nil {
		p.l.Error("unable to serialize serialiproducts", "error", err)
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
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)
	cur := r.URL.Query().Get("currency")

	p.l.Debug("Get record id", id)

	product, err := p.productDB.GetProductByID(id, cur)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Info("[ERROR] fetching producr", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Info("[ERROR] fetching producr", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(product, rw)
	if err != nil {
		p.l.Info("[ERROR] fetching producr", err)

		http.Error(rw, "Failed to marshal product", http.StatusBadRequest)
		return
	}
}
