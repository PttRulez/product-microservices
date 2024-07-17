package main

import (
	"net"
	"os"

	"github.com/pttrulez/product-microservices/currency/data"
	"github.com/pttrulez/product-microservices/currency/protos"
	"github.com/pttrulez/product-microservices/currency/server"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	// create a new gRPC server, use WithInstance to allow http connections
	gs := grpc.NewServer()

	// create an instance of the Currency server
	cs := server.NewCurrency(rates, log)

	// register the currency server
	protos.RegisterCurrencyServer(gs, cs)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connection
	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Failed net")
		os.Exit(1)
	}

	gs.Serve(l)
}
