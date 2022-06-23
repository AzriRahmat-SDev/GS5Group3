package functions

import (
	"GS5Group3/api"
	"fmt"
	"net/http"
)

func Homepage(res http.ResponseWriter, req *http.Request) {

	if alreadyLoggedIn(res, req) {

		fmt.Println("hello i am in homepage")
		if req.Method == http.MethodGet {
			db := api.OpenVenueDB()
			defer db.Close()
			results, err := db.Query("SELECT * FROM plots ")
			if err != nil {
				panic("Error in results")
			}
			defer results.Close()
			var venueArr []api.Plot

			for results.Next() {
				var venue api.Plot
				err := results.Scan(&venue.PlotID, &venue.VenueName, &venue.Address)
				if err != nil {
					panic("Error in scanning")
				}
				venueArr = append(venueArr, venue)
			}
			fmt.Println("hello i am get all plots")
			tpl.ExecuteTemplate(res, "homepage.html", venueArr)
		}

	}
}
