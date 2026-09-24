package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/diffpal/lintpal-demo/internal/orders"
)

func main() {
	apiKey := os.Getenv("DEMO_API_KEY")
	if apiKey == "" {
		log.Fatal("DEMO_API_KEY is required")
	}

	handler := orders.NewHandler(apiKey, orders.NewMemoryStore())
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("orders demo listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
