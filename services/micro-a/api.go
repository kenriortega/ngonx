package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "Port to serve")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Status: %v\n", port)
	})
	r.HandleFunc("/api/v1/version/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Version: %v\n", "0.0.1")
	})

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
