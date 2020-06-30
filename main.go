package main

import (
	"context"
	protos "currency/protos/currency"
	"handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"data"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		//err here
		panic(err)
	}
	defer conn.Close()

	// create currency client
	cc := protos.NewCurrencyClient(conn)

	//create the handler
	ph := handlers.NewProducts(l, v, cc)

	//servermux
	sm := mux.NewRouter()

	//Routes
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	getRouter.Handle("/docs", sh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Starting Server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminte, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(tc)
}
