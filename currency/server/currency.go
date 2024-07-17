package server

import (
	"context"

	"github.com/pttrulez/product-microservices/currency/data"
	"github.com/pttrulez/product-microservices/currency/protos"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	protos.UnimplementedCurrencyServer
	log   hclog.Logger
	rates *data.ExchangeRates
}

func NewCurrency(rates *data.ExchangeRates, log hclog.Logger) *Currency {
	return &Currency{
		log:   log,
		rates: rates,
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}
