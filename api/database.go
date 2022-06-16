package main

import (
	"database/sql"
)

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
