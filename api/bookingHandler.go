package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// booking struct corresponds to the bookings table in the database
type booking struct {
	BookingID      string `json:"BookingID"`
	PlotID         string `json:"PlotID"`
	Username       string `json:"Username"`
	StartDate      string `json:"StartDate"`
	EndDate        string `json:"EndDate"`
	LeaseCompleted string `json:"LeaseCompleted"`
}

const connection string = "root:password@tcp(localhost:32769)/database"

// getHandler handles three types of requests depending on which parameters get passed into the URL.
// If a Username is passed into the URL, this function will pull all bookings associated with that Username.
// If a PlotID is passed into the URL, this function will pull all bookings associated with that PlotID.
// If no paramters are passed into the URL, this function will pull all bookings in the database.
func getHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userParam := params["Username"]
	plotParam := params["PlotID"]

	var query string

	if userParam != "" {
		query = fmt.Sprintf("SELECT * FROM database.bookings WHERE Username='%s'", userParam)
	} else if plotParam != "" {
		query = fmt.Sprintf("SELECT * FROM database.bookings WHERE PlotID='%s'", plotParam)
	} else {
		query = fmt.Sprintf("SELECT * FROM database.bookings")
	}

	bookings := getBookings(query)

	json.NewEncoder(w).Encode(bookings)
}

// getBookings is a support function that takes a SQL query as argument and returns all resulting bookings.
func getBookings(query string) (bookings map[string][]booking) {
	// establish connection to database
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	// create variable to store all bookings
	tempBookings := map[string][]booking{
		"bookings": {},
	}

	for results.Next() {
		var booking booking
		err = results.Scan(&booking.BookingID, &booking.PlotID, &booking.Username, &booking.StartDate, &booking.EndDate, &booking.LeaseCompleted)
		if err != nil {
			panic(err.Error())
		}

		tempBookings["bookings"] = append(tempBookings["bookings"], booking)
	}

	return tempBookings
}

// bookingHandler handles DELETE, POST, PUT, PATCH, and GET requests.
// For GET, it only handles requests for one BookingID.
func bookingHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingParam := params["BookingID"]

	if r.Method == "GET" {
		if bookingExists(bookingParam) {
			query := "SELECT * FROM database.bookings WHERE BookingID = '" + bookingParam + "' LIMIT 1"
			bookings := getBookings(query)
			json.NewEncoder(w).Encode(bookings)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Booking not found"))
		}
	}

	if r.Method == "DELETE" {
		if bookingExists(bookingParam) {
			query := "SELECT * FROM database.bookings WHERE BookingID = '" + bookingParam + "' LIMIT 1"
			bookings := getBookings(query)

			// check if booking has not yet been completed
			leaseCompleted, err := strconv.ParseBool(bookings["bookings"][0].LeaseCompleted)
			if err != nil {
				panic(err.Error())
			}

			if !leaseCompleted {
				// establish connection to database
				db, err := sql.Open("mysql", connection)
				if err != nil {
					panic(err.Error())
				}
				defer db.Close()

				query := fmt.Sprintf("DELETE FROM database.bookings WHERE BookingID='%s'", bookingParam)

				_, err = db.Query(query)
				if err != nil {
					panic(err.Error())
				}

				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("202 - Booking canceled: " + bookingParam))
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Booking has already been completed"))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Booking does not exist"))
		}
	}

	if r.Method == "PATCH" {
		if bookingExists(bookingParam) {
			query := "SELECT * FROM database.bookings WHERE BookingID = '" + bookingParam + "' LIMIT 1"
			bookings := getBookings(query)

			// check if booking has not yet been completed
			leaseCompleted, err := strconv.ParseBool(bookings["bookings"][0].LeaseCompleted)
			if err != nil {
				panic(err.Error())
			}

			if !leaseCompleted {
				// establish connection to database
				db, err := sql.Open("mysql", connection)
				if err != nil {
					panic(err.Error())
				}
				defer db.Close()

				query := fmt.Sprintf("UPDATE database.bookings SET LeaseCompleted='true' WHERE BookingID='%s'", bookingParam)

				_, err = db.Query(query)
				if err != nil {
					panic(err.Error())
				}

				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("202 - Booking " + bookingParam + " has been successfully marked as completed"))
			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Booking has already been completed"))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Booking does not exist"))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {

		if r.Method == "POST" {
			// read data received from client
			var newBooking booking
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newBooking)

				// check if all fields have been received
				if newBooking.PlotID == "" || newBooking.Username == "" || newBooking.StartDate == "" || newBooking.EndDate == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - All fields must be filled out and in JSON format"))
					return
				}

				// check if plot is available on desired dates
				if startDateIsBeforeEndDate(newBooking.StartDate, newBooking.EndDate) && plotAvailable(newBooking.PlotID, newBooking.StartDate, newBooking.EndDate, "") {
					// establish connection to database
					db, err := sql.Open("mysql", connection)
					if err != nil {
						panic(err.Error())
					}
					defer db.Close()

					query := fmt.Sprintf("INSERT INTO bookings (PlotID, Username, StartDate, EndDate, LeaseCompleted) VALUES ('%s', '%s', '%s', '%s', 'false')", newBooking.PlotID, newBooking.Username, newBooking.StartDate, newBooking.EndDate)
					_, err = db.Query(query)
					if err != nil {
						panic(err.Error())
					}

					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Booking added on plot " + newBooking.PlotID))
				} else {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Desired dates are not available or start date is not before end date"))
				}

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply information in JSON format"))
			}
		}

		if r.Method == "PUT" {
			// read data received from client
			var editBooking booking
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &editBooking)

				// check if all fields have been received
				if editBooking.BookingID == "" || editBooking.PlotID == "" || editBooking.Username == "" || editBooking.StartDate == "" || editBooking.EndDate == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - All fields must be filled out and in JSON format"))
					return
				}

				// check if plot is available for new dates
				if startDateIsBeforeEndDate(editBooking.StartDate, editBooking.EndDate) && plotAvailable(editBooking.PlotID, editBooking.StartDate, editBooking.EndDate, editBooking.BookingID) {
					// establish connection to database
					db, err := sql.Open("mysql", connection)
					if err != nil {
						panic(err.Error())
					}
					defer db.Close()

					query := fmt.Sprintf("UPDATE database.bookings SET StartDate='%s', EndDate='%s' WHERE BookingID='%s'", editBooking.StartDate, editBooking.EndDate, editBooking.BookingID)
					_, err = db.Query(query)
					if err != nil {
						panic(err.Error())
					}

					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Booking dates updated for plot " + editBooking.PlotID))
				} else {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Desired dates are not available or start date is not before end date"))
				}

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply information in JSON format"))
			}
		}
	}
}

// bookingExists takes a BookingID and checks whether it currently exists in the database.
func bookingExists(booking string) (exists bool) {
	// establish connection to database
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	search, err := db.Query("SELECT EXISTS(SELECT * FROM database.bookings WHERE BookingID = '" + booking + "')")
	if err != nil {
		panic(err.Error())
	}

	for search.Next() {
		err = search.Scan(&exists)
		if err != nil {
			panic(err.Error())
		}
	}

	return exists
}

// plotAvailable takes PlotID and desired start and end dates and returns a bool noting whether the desired dates are available for booking.
func plotAvailable(plotID, startDate, endDate, bookingID string) (available bool) {
	startDateDate, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Println(err)
	}
	endDateDate, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Println(err)
	}

	query := fmt.Sprintf("SELECT * FROM database.bookings WHERE PlotID='%s' AND LeaseCompleted='false'", plotID)
	bookings := getBookings(query)

	// exclude your own bookingID from unavailable date ranges
	if bookingID != "" {
		tempBookings := map[string][]booking{
			"bookings": {},
		}

		for _, v := range bookings["bookings"] {
			if v.BookingID != bookingID {
				tempBookings["bookings"] = append(tempBookings["bookings"], v)
			}
		}
		bookings = tempBookings
	}

	available = true
	for _, v := range bookings["bookings"] {
		tempStartDate, err := time.Parse("2006-01-02", v.StartDate)
		if err != nil {
			fmt.Println(err)
		}
		tempEndDate, err := time.Parse("2006-01-02", v.EndDate)
		if err != nil {
			fmt.Println(err)
		}

		if startDateDate.Before(tempEndDate) && startDateDate.After(tempStartDate) {
			available = false
			break
		}

		if endDateDate.Before(tempEndDate) && endDateDate.After(tempStartDate) {
			available = false
			break
		}
	}

	return available
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
