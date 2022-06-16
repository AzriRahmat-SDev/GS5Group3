package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Plot struct {
	PlotID    string `json:"PlotID"`
	VenueInfo VenueInformation
}
type VenueInformation struct {
	VenueName string `json:"VenueName"`
	Address   string `json:"Address"`
}

var plotMap map[string]VenueInformation

var PlotList []Plot

func OpenVenueDB() *sql.DB {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/venue_db")

	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Print("Connected")
		// log open
	}
	return db
}

func InsertPlot(db *sql.DB, p Plot) {
	query := fmt.Sprintf("INSERT INTO plots (PlotID, VenueName, Address) VALUES ('%s', '%s', '%s')", p.PlotID, p.VenueInfo.VenueName, p.VenueInfo.Address)
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

func FillMap(db *sql.DB) {
	for k := range plotMap {
		delete(plotMap, k)
	}
	results, err := db.Query("SELECT * FROM plots")
	if err != nil {
		fmt.Println(err.Error())
	}

	for results.Next() {
		var p Plot
		err := results.Scan(&p.PlotID, &p.VenueInfo.VenueName, &p.VenueInfo.Address)
		if err != nil {
			fmt.Println(err.Error())
		}
		plotMap[p.PlotID] = p.VenueInfo
	}
}

func main() {

	plotMap = make(map[string]VenueInformation)
	RunTests()
}

func RunTests() {
	s := Plot{
		PlotID: "ALJ027",
		VenueInfo: VenueInformation{
			VenueName: "Aljunied Park",
			Address:   "Aljunied Road, Happy Garden Estate, 389842",
		},
	}

	InsertPlot(OpenVenueDB(), s)
	EditPlotAddress(OpenVenueDB(), "ALJ001", "Aljunieeed Road, Happy Garden Estate, 389842")
	EditPlotAddress(OpenVenueDB(), "ALJ001", "Aljunied Road, Happy Garden Estate, 389842")
	EditPlotVenueName(OpenVenueDB(), "ALJ001", "Aljunieeed")
	EditPlotVenueName(OpenVenueDB(), "ALJ001", "Aljunied Park")
	DeletePlot(OpenVenueDB(), "ALJ027")
	//RefreshPlots()
}

func RefreshPlots() {

	FillMap(OpenVenueDB())

	for k := range plotMap {
		p := Plot{
			PlotID: k,
			VenueInfo: VenueInformation{
				VenueName: plotMap[k].VenueName,
				Address:   plotMap[k].Address,
			},
		}
		PlotList = append(PlotList, p)
	}
	for x, y := range PlotList {
		fmt.Println(x, y)
	}
}
