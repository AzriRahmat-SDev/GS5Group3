package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const baseURL = "http://localhost:5000/api/v1/"

func startServer() {

	//Allocating memory space to map
	plotMap = make(map[string]Plot)
	//RunTests()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/plots", GetAllPlots).Methods("GET")
	// router.HandleFunc("/api/v1/plots/{plotid}/{address}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")
	// router.HandleFunc("/api/v1/plots/{plotid}/{venuename}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/plots/{plotid}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")

	router.HandleFunc("/api/v1/bookings", getBookings).Methods("GET")
	router.HandleFunc("/api/v1/bookings/{BookingID}", bookingHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/test", Test).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))

}

func main() {
	startServer()

}
