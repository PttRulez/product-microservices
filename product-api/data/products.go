package data

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/pttrulez/product-microservices/currency/protos"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for this user
	//
	// required: true
	// min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

// Products defines a slice of Product
type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	return &ProductsDB{c, l}
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// GetProducts return all products from the database
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	// get exchange rate
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies_USD,
	}
	resp, err := p.currency.GetRate(context.Background(), rr)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * resp.Rate
		pr = append(pr, &np)
	}
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := GetProductById(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

var ErrProductNotFound = errors.New("Not found")

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

// GetProductByID return a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func GetProductById(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = Products{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc23",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
