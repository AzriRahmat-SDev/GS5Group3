package functions

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type SpaceInfo struct {
	VenueName string
	PlotID    string
	Address   string
}

var Space map[string]SpaceInfo

func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("//name of DB", "")
	return db, err
}

func PopulateData(db *sql.DB) {
	results, err := db.Query("SELECT * FROM PlotID")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var space SpaceInfo
		err = results.Scan(&space.VenueName, &space.PlotID, &space.Address)
		if err != nil {
			panic(err.Error())
		}
		Space[space.PlotID] = space
	}

}

func GetAllSpaces(w http.ResponseWriter, r *http.Request) {

	db, err := OpenDatabase()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	PopulateData(db)
	json.NewEncoder(w).Encode(Space)

}
