package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/plots", getPlots).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))
}

func getPlots(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
