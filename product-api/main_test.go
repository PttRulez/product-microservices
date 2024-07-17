package main

import (
	"fmt"
	"testing"

	"github.com/pttrulez/product-microservices/product_api/sdk/client/products"

	"github.com/pttrulez/product-microservices/product_api/sdk/client"
)

func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prods, err := c.Products.ListProducts(params)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("NICE", prods)
}
