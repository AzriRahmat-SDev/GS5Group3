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
}

var plotList []string

var VenueInformationList []VenueInformation

func ViewVenues(res http.ResponseWriter, req *http.Request) {
	// shows all at first
	// if there's a filter, do the filter before changing the range.
	fillVenuesList()

	tpl.ExecuteTemplate(res, "venues.html", VenueInformationList)
}

func fillVenuesList() {
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

	venueName := req.FormValue("VenueName")
	plotList := getPlotList(venueName)
	var venuePlotInfoList []VenuePlotInfo
	sortPlotIDs(plotList, 0, len(plotList)-1)

	for _, k := range plotList {
		x := VenuePlotInfo{
			VenueName: venueName,
			PlotID:    k,
		}
		venuePlotInfoList = append(venuePlotInfoList, x)
	}
	fmt.Println(venuePlotInfoList)
	tpl.ExecuteTemplate(res, "viewvenueplots.html", venuePlotInfoList)
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
