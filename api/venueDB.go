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

var plotMap map[string]Plot
var plotList []Plot

/* Venue Name + Address
 */
var venueMap map[string]string

var initialized bool = false

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

// wip dont touch yet
func nextPlotID(db *sql.DB, venue string) int {
	return 0
}

func populateData(db *sql.DB) {
	for k := range plotMap {
		delete(plotMap, k)
	}
	results, err := db.Query("SELECT * FROM plots")
	if err != nil {
		fmt.Println(err.Error())
	}

	for results.Next() {
		var p Plot
		err := results.Scan(&p.PlotID, &p.VenueName, &p.Address)
		if err != nil {
			fmt.Println(err.Error())
		}
		plotMap[p.PlotID] = p
		plotList = append(plotList, p)
		if _, ok := venueMap[p.VenueName]; !ok {
			venueMap[p.VenueName] = p.Address
		}
	}
	// exposing for template usage
	// for v, k := range venueMap {
	// 	vi := &VenueInformation{VenueName: v, Address: k}
	// 	VenueInformationList = append(VenueInformationList, *vi)
	// }
	// sortVenueInfoList(VenueInformationList, 0, len(VenueInformationList)-1)
}

// func GetVenueInformationList() []VenueInformation {
// 	return VenueInformationList
// }
