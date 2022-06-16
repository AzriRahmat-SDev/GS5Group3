package api

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

	log.Fatal(http.ListenAndServe(":5000", router))
}
