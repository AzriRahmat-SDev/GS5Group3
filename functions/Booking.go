package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type allInfo struct {
	Username      string
	Name          string
	Email         string
	PlotID        string
	VenueName     string
	Address       string
	BookingID     string
	StartDate     string
	EndDate       string
	CurrentLeases bookings
}

type booking struct {
	BookingID      string `json:"BookingID"`
	PlotID         string `json:"PlotID"`
	Username       string `json:"Username"`
	StartDate      string `json:"StartDate"`
	EndDate        string `json:"EndDate"`
	LeaseCompleted string `json:"LeaseCompleted"`
}

type bookings struct {
	Bookings []booking
}

type Plot struct {
	PlotID    string `json:"PlotID"`
	VenueName string `json:"VenueName"`
	Address   string `json:"Address"`
}

const apiURL string = "http://localhost:5000/api/v1/"
const baseURL string = "http://localhost:5000/api/v1/bookings/"
const plotsAPI string = "http://localhost:5000/api/v1/plots/"

func NewBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	Username := req.FormValue("user")
	PlotID := req.FormValue("plot")

	// pull bookings for plot
	leases := callBookingsAPI("byPlot", PlotID)
	currentLeases := onlyCurrentLeases(leases)

	// pull plot info
	plot := callPlotsAPI(PlotID)

	allInfo := map[string]allInfo{
		"allInfo": {
			Username:      Username,
			Name:          Username,
			Email:         Username,
			PlotID:        plot.PlotID,
			VenueName:     plot.VenueName,
			Address:       plot.Address,
			StartDate:     "",
			EndDate:       "",
			CurrentLeases: currentLeases,
		},
	}

	// create new booking when form is submitted
	if req.Method == http.MethodPost {
		StartDate := req.FormValue("StartDate")
		EndDate := req.FormValue("EndDate")

		jsonBooking := packageBookingJSON("", PlotID, Username, StartDate, EndDate)

		response, err := http.Post(apiURL+"bookings/booking/all", "application/json", jsonBooking)

		if err == nil {
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(data))

			if response.StatusCode == 201 {
				fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
			} else {
				fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
				return
			}

		} else {
			fmt.Println(err)
		}

		response.Body.Close()
	}

	tpl.ExecuteTemplate(res, "newbooking.gohtml", allInfo)
}

func EditBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	BookingID := req.FormValue("booking")
	Username := req.FormValue("user")

	// pull bookings for plot
	currentBooking := callBookingsAPI("byBooking", BookingID)
	leases := callBookingsAPI("byPlot", currentBooking.Bookings[0].PlotID)
	currentLeases := onlyCurrentLeases(leases)

	// pull plot info
	plot := callPlotsAPI(currentBooking.Bookings[0].PlotID)

	allInfo := map[string]allInfo{
		"allInfo": {
			Username:      Username,
			Name:          Username,
			Email:         Username,
			PlotID:        plot.PlotID,
			VenueName:     plot.VenueName,
			Address:       plot.Address,
			BookingID:     BookingID,
			StartDate:     currentBooking.Bookings[0].StartDate,
			EndDate:       currentBooking.Bookings[0].EndDate,
			CurrentLeases: currentLeases,
		},
	}

	// edit booking when form is submitted
	if req.Method == http.MethodPost {
		StartDate := req.FormValue("StartDate")
		EndDate := req.FormValue("EndDate")

		jsonBooking := packageBookingJSON(BookingID, plot.PlotID, Username, StartDate, EndDate)

		request, err := http.NewRequest(http.MethodPut, apiURL+"bookings/booking/"+BookingID, jsonBooking)
		if err != nil {
			fmt.Println(err)
		}

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)

		if err == nil {
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(data))

			if response.StatusCode == 201 {
				http.Redirect(res, req, "/editbooking/?booking="+BookingID, http.StatusSeeOther)
			} else {
				fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
				return
			}

		} else {
			fmt.Println(err)
		}

		response.Body.Close()
	}

	tpl.ExecuteTemplate(res, "editbooking.gohtml", allInfo)
}

func DeleteBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	BookingID := req.FormValue("booking")

	request, err := http.NewRequest(http.MethodDelete, baseURL+"booking/"+BookingID, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(data))

		if response.StatusCode == 202 {
			http.Redirect(res, req, "/homepage/", http.StatusSeeOther)
		} else {
			fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
			return
		}

	} else {
		fmt.Println(err)
	}

	response.Body.Close()
}

func callBookingsAPI(byCriteria, criteria string) (bookings bookings) {
	var response *http.Response

	switch byCriteria {
	case "byPlot":
		res, err := http.Get(apiURL + "bookings/plot/" + criteria)
		if err != nil {
			fmt.Println(err)
		} else {
			response = res
		}
	case "byBooking":
		res, err := http.Get(apiURL + "bookings/booking/" + criteria)
		if err != nil {
			fmt.Println(err)
		} else {
			response = res
		}
	default:
		fmt.Println("callBookingsAPI function was not used correctly.")
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(data, &bookings)

	response.Body.Close()

	return bookings
}

func onlyCurrentLeases(bookings bookings) (currentLeases bookings) {
	for _, v := range bookings.Bookings {
		leaseBool, err := strconv.ParseBool(v.LeaseCompleted)
		if err != nil {
			fmt.Println(err)
		}

		if leaseBool == false {
			currentLeases.Bookings = append(currentLeases.Bookings, v)
		}
	}
	return currentLeases
}

func callPlotsAPI(PlotID string) (plot Plot) {
	response, err := http.Get(apiURL + "plots/" + PlotID)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(data, &plot)

		response.Body.Close()
	} else {
		fmt.Println(err)
	}

	return plot
}

func packageBookingJSON(BookingID, PlotID, Username, StartDate, EndDate string) (jsonBooking *bytes.Buffer) {
	booking := booking{
		BookingID: BookingID,
		PlotID:    PlotID,
		Username:  Username,
		StartDate: StartDate,
		EndDate:   EndDate,
	}

	byteBooking, err := json.Marshal(booking)
	if err != nil {
		fmt.Println(err)
	}

	jsonBooking = bytes.NewBuffer(byteBooking)
	return jsonBooking
}
