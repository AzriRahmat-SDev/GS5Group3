package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const plotsAPI string = "http://localhost:5001/api/v1/plots/"

type VenueInformation struct {
	VenueName string `json:"VenueName"`
	Address   string `json:"Address"`
}

var venueMap map[string]string

type VenuePlotInfo struct {
	VenueName string
	PlotID    string
	Bookings  []booking
}

var plotList []string

var VenueInformationList []VenueInformation

// var venueListPulled bool = false

func ViewVenues(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(res, req) {
		// if !venueListPulled {
		// to be called again by admins if any changes are made.
		fillVenuesList()
		// 	venueListPulled = true
		// }
		tpl.ExecuteTemplate(res, "venues.html", VenueInformationList)
	}
}
func fillVenuesList() {

	venueMap = make(map[string]string)
	VenueInformationList = []VenueInformation{}
	res, err := http.Get(plotsAPI + "venue/")
	if err == nil {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(data, &venueMap)

		res.Body.Close()
	} else {
		fmt.Println(err)
	}
	for v, k := range venueMap {
		vi := &VenueInformation{VenueName: v, Address: k}
		VenueInformationList = append(VenueInformationList, *vi)
	}
	sortVenueInfoList(VenueInformationList, 0, len(VenueInformationList)-1)
}

func sortVenueInfoList(v []VenueInformation, left int, right int) {
	if left < right {
		pivotIndex := venueInfoPartition(v, left, right)
		sortVenueInfoList(v, left, pivotIndex-1)
		sortVenueInfoList(v, pivotIndex+1, right)
	}
}

func venueInfoPartition(v []VenueInformation, left int, right int) int {
	pivot := strings.ToLower(v[(left+right)/2].VenueName)
	for left < right {
		for strings.ToLower(v[left].VenueName) < pivot {
			left++
		}
		for strings.ToLower(v[right].VenueName) > pivot {
			right--
		}
		if left < right {
			temp := v[left]
			v[left] = v[right]
			v[right] = temp
		}
	}
	return left
}

// Viewing Specific Venue with PlotID Information:
func ViewVenuePlots(res http.ResponseWriter, req *http.Request) {

	if alreadyLoggedIn(res, req) {
		venueName := req.FormValue("VenueName")
		plotList := getPlotList(venueName)
		var venuePlotInfoList []VenuePlotInfo
		sortPlotIDs(plotList, 0, len(plotList)-1)

		for _, k := range plotList {
			bookings := callBookingsAPI("byPlot", k)
			bookings = onlyCurrentLeases(bookings)

			x := VenuePlotInfo{
				VenueName: venueName,
				PlotID:    k,
				Bookings:  bookings.Bookings,
			}
			venuePlotInfoList = append(venuePlotInfoList, x)
		}

		venuePlotInfoListFinal := venuePlotInfoList

		// filter plots by date when form is submitted
		if req.Method == http.MethodPost {
			StartDate := req.FormValue("StartDate")
			EndDate := req.FormValue("EndDate")

			if StartDate != "" && EndDate != "" {
				if !startDateIsBeforeEndDate(StartDate, EndDate) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Start date is not before end date"))
					return
				}

				var tempPlots []VenuePlotInfo
				for _, v := range venuePlotInfoList {

					available := true

					for _, value := range v.Bookings {
						if !plotAvailable(StartDate, EndDate, value.StartDate, value.EndDate) {
							available = false
							break
						}
					}

					if available == true {
						tempPlots = append(tempPlots, v)
					}

				}
				venuePlotInfoListFinal = tempPlots
			}

		}

		tpl.ExecuteTemplate(res, "viewvenueplots.html", venuePlotInfoListFinal)
	}
}

func getPlotList(venueName string) []string {
	var r []string
	req, err := http.Get(plotsAPI + "venue/" + venueName)
	if err != nil {
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {

	}
	json.Unmarshal(data, &r)
	req.Body.Close()
	return r
}

//sort
func sortPlotIDs(v []string, left int, right int) {
	if left < right {
		pivotIndex := plotIDPartitions(v, left, right)
		sortPlotIDs(v, left, pivotIndex-1)
		sortPlotIDs(v, pivotIndex+1, right)
	}
}

func plotIDPartitions(v []string, left int, right int) int {
	pivot := strings.ToLower(v[(left+right)/2])
	for left < right {
		for strings.ToLower(v[left]) < pivot {
			left++
		}
		for strings.ToLower(v[right]) > pivot {
			right--
		}
		if left < right {
			temp := v[left]
			v[left] = v[right]
			v[right] = temp
		}
	}
	return left
}
