package main

import (
	"log"
	"net/http"

	"samplecalculatorproject/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/capabilities", api.Capabilities)
	mux.HandleFunc("/add", api.Add)
	mux.HandleFunc("/subtract", api.Subtract)
	mux.HandleFunc("/multiply", api.Multiply)
	mux.HandleFunc("/divide", api.Divide)

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
