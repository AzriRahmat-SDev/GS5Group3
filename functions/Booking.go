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

const baseURL string = "http://localhost:5000/api/v1/bookings/"
const plotsAPI string = "http://localhost:5000/api/v1/plots/"

func NewBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	Username := req.FormValue("user")
	PlotID := req.FormValue("plot")

	// pull bookings for plot
	var tempLeases bookings

	response, err := http.Get(baseURL + "plot/" + PlotID)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(data, &tempLeases)

		response.Body.Close()
	} else {
		fmt.Println(err)
	}

	var availableLeases bookings

	for _, v := range tempLeases.Bookings {
		leaseBool, err := strconv.ParseBool(v.LeaseCompleted)
		if err != nil {
			fmt.Println(err)
		}

		if leaseBool == false {
			availableLeases.Bookings = append(availableLeases.Bookings, v)
		}
	}

	// pull plot info
	var plot Plot

	response, err = http.Get(plotsAPI + PlotID)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		data = data[3:] // remove method that gets returned in front of JSON object
		json.Unmarshal(data, &plot)

		response.Body.Close()
	} else {
		fmt.Println(err)
	}

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
			CurrentLeases: availableLeases,
		},
	}

	// create new booking when form is submitted
	if req.Method == http.MethodPost {
		StartDate := req.FormValue("StartDate")
		EndDate := req.FormValue("EndDate")

		newBooking := booking{
			PlotID:    PlotID,
			Username:  Username,
			StartDate: StartDate,
			EndDate:   EndDate,
		}

		byteBooking, err := json.Marshal(newBooking)
		if err != nil {
			fmt.Println(err)
		}
		jsonBooking := bytes.NewBuffer(byteBooking)

		response, err := http.Post(baseURL+"booking/all", "application/json", jsonBooking)

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
	PlotID := req.FormValue("plot")

	var tempLeases bookings

	response, err := http.Get(baseURL + "plot/" + PlotID)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(data, &tempLeases)

		response.Body.Close()
	} else {
		fmt.Println(err)
	}

	var availableLeases bookings

	for _, v := range tempLeases.Bookings {
		leaseBool, err := strconv.ParseBool(v.LeaseCompleted)
		if err != nil {
			fmt.Println(err)
		}

		if leaseBool == false {
			availableLeases.Bookings = append(availableLeases.Bookings, v)
		}
	}

	allInfo := map[string]allInfo{
		"allInfo": {
			Username:      Username,
			Name:          Username,
			Email:         Username,
			PlotID:        PlotID,
			VenueName:     PlotID,
			Address:       PlotID,
			BookingID:     BookingID,
			StartDate:     "",
			EndDate:       "",
			CurrentLeases: availableLeases,
		},
	}

	// edit booking when form is submitted
	if req.Method == http.MethodPost {
		StartDate := req.FormValue("StartDate")
		EndDate := req.FormValue("EndDate")

		editBooking := booking{
			BookingID: BookingID,
			PlotID:    PlotID,
			Username:  Username,
			StartDate: StartDate,
			EndDate:   EndDate,
		}

		byteBooking, err := json.Marshal(editBooking)
		if err != nil {
			fmt.Println(err)
		}
		jsonBooking := bytes.NewBuffer(byteBooking)

		request, err := http.NewRequest(http.MethodPut, baseURL+"booking/"+BookingID, jsonBooking)
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
				fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
				return
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
	Username := req.FormValue("user")
	PlotID := req.FormValue("plot")

	var tempLeases bookings

	response, err := http.Get(baseURL + "plot/" + PlotID)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(data, &tempLeases)

		response.Body.Close()
	} else {
		fmt.Println(err)
	}

	var availableLeases bookings

	for _, v := range tempLeases.Bookings {
		leaseBool, err := strconv.ParseBool(v.LeaseCompleted)
		if err != nil {
			fmt.Println(err)
		}

		if leaseBool == false {
			availableLeases.Bookings = append(availableLeases.Bookings, v)
		}
	}

	allInfo := map[string]allInfo{
		"allInfo": {
			Username:      Username,
			Name:          Username,
			Email:         Username,
			PlotID:        PlotID,
			VenueName:     PlotID,
			Address:       PlotID,
			BookingID:     BookingID,
			StartDate:     "",
			EndDate:       "",
			CurrentLeases: availableLeases,
		},
	}

	// delete booking when form is submitted
	if req.Method == http.MethodPost {

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
				fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
				return
			} else {
				fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
				return
			}

		} else {
			fmt.Println(err)
		}

		response.Body.Close()
	}

	tpl.ExecuteTemplate(res, "deletebooking.gohtml", allInfo)
}
