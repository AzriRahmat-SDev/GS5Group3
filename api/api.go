package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//Allocating memory space to map
	plotMap = make(map[string]VenueInformation)
	RunTests()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/plots", GetAllPlots).Methods("GET")
	router.HandleFunc("/api/v1/plots/{plotid}/{address}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/plots/{plotid}/{VenueName}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/bookings", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/booking/{BookingID}", bookingHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/bookings/user/{UserID}", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/plot/{PlotID}", getHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))
}
