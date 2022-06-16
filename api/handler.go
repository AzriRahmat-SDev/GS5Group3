package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

type SpaceInfo struct {
	VenueName string
	PlotID    string
	Address   string
}

//Declaring Map Existence globally
var Space map[string]SpaceInfo

func GetAllPlots(w http.ResponseWriter, r *http.Request) {

	db, err := OpenDatabase()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	PopulateData(db)
	json.NewEncoder(w).Encode(Space)

}

func PlotHandler(w http.ResponseWriter, r *http.Request) {

	db, err := OpenDatabase()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	params := mux.Vars(r)

	if r.Method == "GET" {
		PopulateData(db)
		if _, ok := Space[params["plotid"]]; ok {
			json.NewEncoder(w).Encode(Space[params["plotid"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found from GET"))
		}
	}

	if r.Method == "DELETE" {
		PopulateData(db)
		if key, ok := Space[params["plotid"]]; ok {
			query := fmt.Sprintf("DELETE FROM plot WHERE ID='%s'", key.PlotID)
			_, err := db.Query(query)
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("Delete is successful")
			delete(Space, params["plotid"])
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Plot deleted: " + params["plotid"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	}
}
