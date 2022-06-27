// Package api is an API server that runs independent of any client.
package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const baseURL = "http://localhost:5001/api/v1/"

func StartServer() {

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/users/", getAllUsers).Methods("GET")
	router.HandleFunc("/api/v1/users/update/{username}", userHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/plots/", getAllPlots).Methods("GET")
	router.HandleFunc("/api/v1/plots/{plotid}", plotHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/plots/venue/", venueHandler).Methods("GET")
	router.HandleFunc("/api/v1/plots/venue/{VenueName}", viewVenuePlots).Methods("GET")

	router.HandleFunc("/api/v1/bookings", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/user/{Username}", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/plot/{PlotID}", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/booking/{BookingID}", bookingHandler).Methods("GET", "POST", "PUT", "DELETE", "PATCH")

	log.Fatal(http.ListenAndServe(":5001", router))

}
