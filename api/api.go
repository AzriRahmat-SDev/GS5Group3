package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const baseURL = "http://localhost:5000/api/v1/"

func startServer() {

	//Allocating memory space to map
	plotMap = make(map[string]VenueInformation)
	//RunTests()

	router := mux.NewRouter()
  router.HandleFunc("/api/v1/plots", GetAllPlots).Methods("GET")
	// router.HandleFunc("/api/v1/plots/{plotid}/{address}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")
	// router.HandleFunc("/api/v1/plots/{plotid}/{venuename}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/plots/{plotid}", PlotHandler).Methods("GET", "POST", "PUT", "DELETE")

	router.HandleFunc("/api/v1/bookings", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/booking/{BookingID}", bookingHandler).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/api/v1/bookings/user/{UserID}", getHandler).Methods("GET")
	router.HandleFunc("/api/v1/bookings/plot/{PlotID}", getHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))

}

func main() {
	startServer()

}

func Test(w http.ResponseWriter, r *http.Request) {
	var p VenueInformation
	PopulateData(OpenVenueDB())
	s := plotMap["ALJ001"].VenueName
	s = strings.ReplaceAll(s, " ", "%20")

	request, _ := http.NewRequest(http.MethodDelete, baseURL+"/plots/ALJ001/", nil)
	client := &http.Client{}

	resp, _ := client.Do(request)
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &p)
	fmt.Println("DATA =", p, s)
	resp.Body.Close()
}
