package main

import (
	"currency/protos"
	"currency/server"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func main() {
	log := hclog.Default()

	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, cs)

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Error("Failed net")
		os.Exit(1)
	}

	gs.Serve(l)
}
