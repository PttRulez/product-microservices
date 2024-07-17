package handlers

import (
	"context"
	"net/http"

	"github.com/pttrulez/product-microservices/currency/protos"
	"github.com/pttrulez/product-microservices/product_api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

// GetProducts return the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("Handle GET Products")
	rw.Header().Add("Content-Type", "application/json")

	// fetch list of products from datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := data.ToJSON(lp, rw)
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
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Info("[DEBUG] get record id", id)

	product, _, err := data.GetProductById(id)

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

	// get exchange rate
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies_USD,
	}
	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Error("[Error] error getting new rate", err)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	product.Price = resp.Rate * product.Price

	err = data.ToJSON(product, rw)
	if err != nil {
		p.l.Info("[ERROR] fetching producr", err)

		http.Error(rw, "Failed to marshal product", http.StatusBadRequest)
		return
	}
}
