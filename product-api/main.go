package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/pttrulez/product-microservices/product_api/data"
	"github.com/pttrulez/product-microservices/product_api/handlers"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.com/pttrulez/product-microservices/currency/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
	// create logger
	l := hclog.Default()

	// create validator
	v := data.NewValidation()

	conn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create client
	cc := protos.NewCurrencyClient(conn)

	// create handlers
	ph := handlers.NewProducts(l, v, cc)

	//create mux
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.Update)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// create server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	l.Info("Running on ")
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
