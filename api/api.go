package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const baseURL = "http://localhost:5000/api/v1/"

func StartServer() {

	//Allocating memory space to map
	plotMap = make(map[string]Plot)
	//RunTests()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/plots", GetAllPlots).Methods("GET")
	router.HandleFunc("/api/v1/plots/{plotid}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")

	router.HandleFunc("/api/v1/bookings", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/booking/{BookingID}", bookingHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/bookings/user/{UserID}", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/plot/{PlotID}", getHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))

}

func main() {
	StartServer()

}
