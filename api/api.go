package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//Allocating memory space to map
	Space = make(map[string]SpaceInfo)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/plots", GetAllPlots).Methods("GET")
	router.HandleFunc("/api/v1/plots/{plotid}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/bookings", getBookings).Methods("GET")
	router.HandleFunc("/api/v1/bookings/{BookingID}", bookingHandler).Methods("GET", "POST", "PUT", "DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
