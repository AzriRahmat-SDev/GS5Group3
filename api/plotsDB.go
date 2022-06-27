package api

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Plot struct {
	PlotID    string `json:"PlotID"`
	VenueName string `json:"VenueName"`
	Address   string `json:"Address"`
}

type PlotMap map[string]Plot

type VenueMap map[string]string

func OpenVenueDB() *sql.DB {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func InsertPlot(db *sql.DB, p Plot) {
	query := fmt.Sprintf("INSERT INTO plots (PlotID, VenueName, Address) VALUES ('%s', '%s', '%s')", p.PlotID, p.VenueName, p.Address)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nSuccessful insert plot @ '%s'", p)
	}
}

func EditPlotAddress(db *sql.DB, plotID string, address string) {
	query := fmt.Sprintf("UPDATE plots SET Address='%s' WHERE PlotID='%s'", address, plotID)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nSuccessful update plot address @ '%s' with '%s'", plotID, address)
	}
}

func EditPlotVenueName(db *sql.DB, plotID string, venueName string) {
	query := fmt.Sprintf("UPDATE plots SET VenueName='%s' WHERE PlotID='%s'", venueName, plotID)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nSuccessful update venue name @ '%s', %s", plotID, venueName)
	}
}

func DeletePlot(db *sql.DB, plotID string) {
	query := fmt.Sprintf("DELETE FROM plots WHERE PlotID='%s'", plotID)
	_, err := db.Query(query)
	if err != nil {
		log.Panic(err.Error())
	} else {
		fmt.Printf("\nSuccessful delete plot @ '%s'", plotID)
	}
}
