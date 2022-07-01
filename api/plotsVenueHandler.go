package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/gorilla/mux"
)

/*
 /api/v1/plots runs GetAllPlots (GET)
 Gets all plots and returns it as a json with value {PlotID:{PlotID, VenueName, Address}}
*/
func GetAllPlots(w http.ResponseWriter, r *http.Request) {

	pMap := makePlotMap("")
	json.NewEncoder(w).Encode(pMap)
}

/*
/api/v1/plots/venue runs venueHandler
Gets all venues and returns it as a json with value {VenueName:Address}
*/
func venueHandler(w http.ResponseWriter, r *http.Request) {
	venueMap := makeVenueMap()
	json.NewEncoder(w).Encode(venueMap)
}

/*
/api/v1/plots/venue/{VenueName} runs viewVenuePlots
GET Venue PlotIDs based on VenueName given, then returns a json of PlotIDs.
*/
func viewVenuePlots(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if plotDBRowExists(params["VenueName"], "VenueName") {
		var plotIDs []string
		pm := makePlotMap("")
		for v, k := range pm {
			if k.VenueName == params["VenueName"] {
				plotIDs = append(plotIDs, v)
			}
		}
		json.NewEncoder(w).Encode(plotIDs)
	}
}

/*


/api/v1/plots/{plotId} runs PlotHandler (GET, DELETE, POST, PUT)
PlotHandler handles GET, DELETE, POST and PUT for individual PlotIDs:
GET returns a json {PlotID:{PlotID, VenueName, Address}} based off what PlotID was given.
DELETE will delete the respective PlotID given to the handler if it exists.
POST will add a new PlotID, VenueName and Address depending on what plotid was given to it.
PUT will either add a new PlotID if it does not exist, or alter the VenueName and Address depending on what the administrator fills it with.
*/
func PlotHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenVenueDB()
	defer db.Close()

	params := mux.Vars(r)

	if r.Method == "GET" {
		if plotDBRowExists(params["plotid"], "PlotID") {
			fmt.Println("PLOT ID : ", params["plotid"])
			p := makePlotMap(params["plotid"])
			json.NewEncoder(w).Encode(p)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No Plot found from GET"))
		}
	}
	if r.Method == "DELETE" {
		if plotDBRowExists(params["plotid"], "PlotID") {
			DeletePlot(db, params["plotid"])
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Plot deleted: " + params["plotid"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No Plot found"))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			var newPlot Plot
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &newPlot)
				if newPlot.PlotID == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply Venue or Address information in JSON format"))
					return
				}
				//basically reads(from client) if plot id does not exist(in DB)
				if !plotDBRowExists(params["plotid"], "PlotID") {
					InsertPlot(db, newPlot)
					fmt.Println("Insert was successful")
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Plot added: " + params["plotid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Duplicate plot ID"))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply Plot information " + "in JSON format"))
			}
		}
		//PUT request here
		if r.Method == "PUT" {
			var newPlot Plot
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {

				json.Unmarshal(reqBody, &newPlot)

				// if both fields are empty
				if newPlot.VenueName == "" && newPlot.Address == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply Plot " + "information " + "in JSON format"))
					return
				}
				if !plotDBRowExists(params["plotid"], "PlotID") {
					InsertPlot(db, newPlot)
					fmt.Println("Insert was successful")
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Plot added: " + params["plotid"]))
				} else {
					if newPlot.Address != "" {
						EditPlotAddress(db, params["plotid"], newPlot.Address)
					}
					if newPlot.VenueName != "" {
						EditPlotVenueName(db, params["plotid"], newPlot.VenueName)
					}
					w.Write([]byte("201-" + params["plotid"] + " has been updated." + newPlot.Address + newPlot.VenueName))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply Plot information " + "in JSON format"))
			}
		}
	}
}

/*
Helper to find if val exists in column.
*/
func plotDBRowExists(val string, column string) bool {
	db := OpenVenueDB()
	defer db.Close()
	r := false
	s, err := db.Query("SELECT EXISTS(SELECT * FROM database.plots WHERE " + column + "='" + val + "')")
	if err != nil {
		panic(err.Error())
	}
	for s.Next() {
		err = s.Scan(&r)
		if err != nil {
			panic(err.Error())
		}
	}
	return r
}

// Returns map values for Plot GET
func makePlotMap(val string) PlotMap {
	plotMap := make(map[string]Plot)
	db := OpenVenueDB()
	defer db.Close()
	if val == "" {
		query := fmt.Sprintf("SELECT * from plots")
		res, err := db.Query(query)
		if err != nil {
		}
		for res.Next() {
			var p Plot
			res.Scan(&p.PlotID, &p.VenueName, &p.Address)
			plotMap[p.PlotID] = p
		}
	} else {
		result, err := db.Query("SELECT * from database.plots WHERE PlotID = '" + val + "'")
		if err != nil {
			fmt.Println(err)
		}
		for result.Next() {
			var p Plot
			err := result.Scan(&p.PlotID, &p.VenueName, &p.Address)
			if err != nil {
				fmt.Println(err.Error())
			}
			plotMap["Plot"] = p
		}
	}
	return plotMap
}

// returns map values for Venue GET
func makeVenueMap() map[string]string {
	venueMap := make(VenueMap)
	db := OpenVenueDB()
	defer db.Close()
	query := fmt.Sprintf("SELECT DISTINCT VenueName, Address from plots")
	res, err := db.Query(query)
	if err != nil {
	}
	for res.Next() {
		var x string
		var y string
		res.Scan(&x, &y)
		venueMap[x] = y

	}
	return venueMap
}
