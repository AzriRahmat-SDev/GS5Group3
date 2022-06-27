// Package functions contains all of the functions that serve pages as well as many support functions.
package functions

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// allInfo struct is a combination of data from all three tables in the database: users, plots, and bookings.
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
	ExpiredLeases bookings
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

type PlotMap map[string]Plot

const apiURL string = "http://localhost:5001/api/v1/"
const connection string = "root:password@tcp(localhost:32769)/database"

// NewBooking serves a page that allows a user to make a new booking.
func NewBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	PlotID := req.FormValue("plot")

	var wg sync.WaitGroup
	wg.Add(3)

	// pull user info
	var user updateUsers
	go func() {
		defer wg.Done()
		cookie, err := req.Cookie("myCookie")
		if err != nil {
			http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
			return
		}
		user = getUser(cookie)
	}()

	// pull bookings for plot
	var currentLeases bookings
	go func() {
		defer wg.Done()
		leases := callBookingsAPI("byPlot", PlotID)
		currentLeases = onlyCurrentLeases(leases)
	}()

	// pull plot info
	var plotMap PlotMap
	go func() {
		defer wg.Done()
		plotMap = callPlotsAPI(PlotID)
	}()

	wg.Wait()

	allInfo := map[string]allInfo{
		"allInfo": {
			Username:      user.Username,
			Name:          user.Name,
			Email:         user.Email,
			PlotID:        plotMap["Plot"].PlotID,
			VenueName:     plotMap["Plot"].VenueName,
			Address:       plotMap["Plot"].Address,
			StartDate:     "",
			EndDate:       "",
			CurrentLeases: currentLeases,
		},
	}

	// create new booking when form is submitted
	if req.Method == http.MethodPost {
		StartDate := req.FormValue("StartDate")
		EndDate := req.FormValue("EndDate")

		jsonBooking := packageBookingJSON("", PlotID, user.Username, StartDate, EndDate)

		response, err := http.Post(apiURL+"bookings/booking/all", "application/json", jsonBooking)

		if err == nil {
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(data))

			if response.StatusCode == 201 {
				http.Redirect(res, req, "/newbooking/?plot="+PlotID, http.StatusSeeOther)
			} else {
				fmt.Fprintf(res, string(data))
				return
			}

		} else {
			fmt.Println(err)
		}

		response.Body.Close()
	}

	tpl.ExecuteTemplate(res, "newbooking.html", allInfo)
}

// EditBooking serves a page that allows the user to edit an existing booking dates.
func EditBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	BookingID := req.FormValue("booking")

	var wg sync.WaitGroup
	wg.Add(1)

	// pull user info
	var user updateUsers
	go func() {
		defer wg.Done()

		cookie, err := req.Cookie("myCookie")
		if err != nil {
			http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
			return
		}
		user = getUser(cookie)

	}()

	// pull bookings for plot
	currentBooking := callBookingsAPI("byBooking", BookingID)
	if len(currentBooking.Bookings) == 0 {
		fmt.Fprintf(res, "Booking does not exist.")
		return
	}

	leases := callBookingsAPI("byPlot", currentBooking.Bookings[0].PlotID)
	currentLeases := onlyCurrentLeases(leases)

	// pull plot info
	plotMap := callPlotsAPI(currentBooking.Bookings[0].PlotID)

	wg.Wait()

	allInfo := map[string]allInfo{
		"allInfo": {
			Username:      user.Username,
			Name:          user.Name,
			Email:         user.Email,
			PlotID:        plotMap["Plot"].PlotID,
			VenueName:     plotMap["Plot"].VenueName,
			Address:       plotMap["Plot"].Address,
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

		jsonBooking := packageBookingJSON(BookingID, plotMap["Plot"].PlotID, user.Username, StartDate, EndDate)

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
				fmt.Fprintf(res, string(data))
				return
			}

		} else {
			fmt.Println(err)
		}

		response.Body.Close()
	}

	tpl.ExecuteTemplate(res, "editbooking.html", allInfo)
}

// DeleteBooking deletes a booking.
func DeleteBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	BookingID := req.FormValue("booking")

	request, err := http.NewRequest(http.MethodDelete, apiURL+"bookings/booking/"+BookingID, nil)
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

// CompleteBooking marks a booking lease as having been completed.
func CompleteBooking(res http.ResponseWriter, req *http.Request) {
	// URL queries
	BookingID := req.FormValue("booking")

	request, err := http.NewRequest(http.MethodPatch, apiURL+"bookings/booking/"+BookingID, nil)
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
			http.Redirect(res, req, "/user/", http.StatusSeeOther)
		} else {
			fmt.Fprintf(res, strconv.Itoa(response.StatusCode))
			return
		}

	} else {
		fmt.Println(err)
	}

	response.Body.Close()
}

// callBookingsAPI allows the client to call the bookings API by PlotID, BookingID, or Username.
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
	case "byUser":
		res, err := http.Get(apiURL + "bookings/user/" + criteria)
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

// onlyCurrentLeases takes a list of bookings and returns only the bookings that have not been marked as completed.
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

// callPlotsAPI allows the client to call the plot API to get information on a single PlotID.
func callPlotsAPI(PlotID string) (PlotMap PlotMap) {
	response, err := http.Get(apiURL + "plots/" + PlotID)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(data, &PlotMap)

		response.Body.Close()
	} else {
		fmt.Println(err)
	}
	fmt.Println(PlotMap)
	return PlotMap
}

// packageBookingJSON takes the information necessary to create or modify a row in the database and packages it into JSON format for the API.
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

// getUser takes the cookie value from the browser and gets user information from the users database.
func getUser(cookie *http.Cookie) (user updateUsers) {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM database.users WHERE Username='%s'", cookie.Value)
	row := db.QueryRow(query)

	err = row.Scan(&user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		fmt.Println(err, "Scan error")
		return
	}

	return user
}

// startDateIsBeforeEndDate takes a start and end date and makes sure the start date is before the end date.
func startDateIsBeforeEndDate(startDate, endDate string) bool {
	startDateDate, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Println(err)
	}
	endDateDate, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Println(err)
	}

	if startDateDate.Before(endDateDate) {
		return true
	} else {
		return false
	}
}

// plotAvailable takes desired start and end dates and returns a bool noting whether the desired dates are available for booking.
func plotAvailable(desiredStartDate, desiredEndDate, bookingStartDate, bookingEndDate string) (available bool) {
	desiredStartDateDate, err := time.Parse("2006-01-02", desiredStartDate)
	if err != nil {
		fmt.Println(err)
	}
	desiredEndDateDate, err := time.Parse("2006-01-02", desiredEndDate)
	if err != nil {
		fmt.Println(err)
	}
	bookingStartDateDate, err := time.Parse("2006-01-02", bookingStartDate)
	if err != nil {
		fmt.Println(err)
	}
	bookingEndDateDate, err := time.Parse("2006-01-02", bookingEndDate)
	if err != nil {
		fmt.Println(err)
	}

	if desiredStartDateDate.Before(bookingEndDateDate) && desiredStartDateDate.After(bookingStartDateDate) {
		return false
	} else if desiredEndDateDate.Before(bookingEndDateDate) && desiredEndDateDate.After(bookingStartDateDate) {
		return false
	} else {
		return true
	}
}
