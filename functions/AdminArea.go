package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Users struct {
	Name     string `field:"Name"`
	Username string `field:"Username"`
	Email    string `field:"Email"`
}

func AllUsers(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedInAdmin(res, req) {
		db := connectUserDB()
		defer db.Close()
		fmt.Println("*****AllUsesHandler running*****")
		if req.Method == http.MethodGet {
			results, err := db.Query("SELECT Name, Username, Email FROM users")
			defer results.Close()
			if err != nil {
				panic("Error in Allusers Query")
			}

			var userArr []Users

			for results.Next() {
				var user Users
				err := results.Scan(&user.Name, &user.Username, &user.Email)
				if err != nil {
					panic("Error in scan")
				}

				userArr = append(userArr, user)
			}

			tpl.ExecuteTemplate(res, "allusers.html", userArr)

		}
	}

}

func DeleteRecord(res http.ResponseWriter, req *http.Request) {

	if alreadyLoggedInAdmin(res, req) {
		db := connectUserDB()
		defer db.Close()
		fmt.Println("*****deleteHandler running*****")
		req.ParseForm()
		username := req.FormValue("username")

		stmt := fmt.Sprintf("DELETE FROM users WHERE (`Username` = '%v')", username)
		result, err := db.Query(stmt)
		defer result.Close()
		if err != nil {
			panic(err)
		}
		InfoLogger.Printf("User %s deleted from user database.", username)
		tpl.ExecuteTemplate(res, "result.html", "User was Successfully Deleted")
	}

}

// If the plotid exists, it will change the venue name and address.
func AddOrEditPlot(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedInAdmin(res, req) {
		if req.Method == http.MethodPost {
			plotID := req.FormValue("plotid")
			venueName := req.FormValue("venuename")
			address := req.FormValue("address")

			p := Plot{
				PlotID:    plotID,
				VenueName: venueName,
				Address:   address,
			}
			r, err := json.Marshal(p)
			input := bytes.NewBuffer(r)
			if err != nil {
				ErrorLogger.Println("Error in json format!")
			}
			request, err := http.NewRequest(http.MethodPut, plotsAPI+p.PlotID, input)
			request.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			response, err := client.Do(request)

			if response.StatusCode == 201 {
				InfoLogger.Println("A method put was made as ", p)
				http.Redirect(res, req, "/allusers", http.StatusSeeOther)
			} else {
				ErrorLogger.Println("Attempted put at ", request.URL, "but failed with error", response.StatusCode)
				return
			}

		}
		tpl.ExecuteTemplate(res, "addeditplot.html", req)
	}
}

func DeletePlot(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedInAdmin(res, req) {
		var plotIDList []string
		fillVenuesList()
		for _, v := range VenueInformationList {
			i := getPlotList(v.VenueName)
			for _, j := range i {
				plotIDList = append(plotIDList, j)
			}
		}
		sortPlotIDs(plotIDList, 0, len(plotIDList)-1)

		if req.Method == http.MethodPost {
			plotToDelete := req.FormValue("plotid")
			request, err := http.NewRequest(http.MethodDelete, plotsAPI+plotToDelete, nil)
			if err != nil {
				ErrorLogger.Println("Attempted delete at ", plotToDelete, "but failed.")
			}
			client := &http.Client{}
			response, err := client.Do(request)
			if response.StatusCode == 202 {
				InfoLogger.Println("Deleted PlotID", plotToDelete)
				http.Redirect(res, req, "/allusers", http.StatusSeeOther)
			} else {
				ErrorLogger.Println("Attempted delete at ", plotToDelete, "but failed with error", response.StatusCode)
				return
			}
		}

		tpl.ExecuteTemplate(res, "deletePlot.html", plotIDList)
	}

}
