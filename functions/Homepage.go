package functions

import (
	"GS5Group3/api"
	"net/http"
)

func Homepage(res http.ResponseWriter, req *http.Request) {

	if alreadyLoggedIn(res, req) {
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
			tpl.ExecuteTemplate(res, "homepage.html", venueArr)
		}

	}
}
