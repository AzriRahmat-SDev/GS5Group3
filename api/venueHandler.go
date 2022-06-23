package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/gorilla/mux"
)

//Declaring Map Existence globally

func GetAllPlots(w http.ResponseWriter, r *http.Request) {
	db := OpenVenueDB()
	defer db.Close()
	PopulateData(db)
	json.NewEncoder(w).Encode(plotMap)
}

func PlotHandler(w http.ResponseWriter, r *http.Request) {

	db := OpenVenueDB()
	defer db.Close()

	params := mux.Vars(r)
	fmt.Fprintf(w, r.Method)

	if r.Method == "GET" {
		PopulateData(db)
		if _, ok := plotMap[params["plotid"]]; ok {
			fmt.Println("PLOT ID : ", params["plotid"])
			json.NewEncoder(w).Encode(plotMap[params["plotid"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No Plot found from GET"))
		}
	}

	if r.Method == "DELETE" {
		PopulateData(db)
		if _, ok := plotMap[params["plotid"]]; ok {
			//DeletePlot(db, params["plotid"])
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Plot deleted: " + params["plotid"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No Plot found"))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		PopulateData(db)
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
				if _, ok := plotMap[params["plotid"]]; !ok {
					InsertPlot(db, newPlot)
					fmt.Println("Insert was successful")
					plotMap[params["plotid"]] = newPlot
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
			PopulateData(db)
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
				if _, ok := plotMap[params["plotid"]]; !ok {
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
