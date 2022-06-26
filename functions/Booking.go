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

	_ "github.com/go-sql-driver/mysql"
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
				http.Redirect(res, req, "/homepage/", http.StatusSeeOther)
			} else {
				fmt.Fprintf(res, string(data))
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

	tpl.ExecuteTemplate(res, "editbooking.gohtml", allInfo)
}

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
