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

type booking struct {
	BookingID      string `json:"BookingID"`
	PlotID         string `json:"PlotID"`
	UserID         string `json:"UserID"`
	StartDate      string `json:"StartDate"`
	EndDate        string `json:"EndDate"`
	LeaseCompleted string `json:"LeaseCompleted"`
}

const connection string = "root:password@tcp(localhost:32769)/database"

func getHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userParam := params["UserID"]
	plotParam := params["PlotID"]

	var query string

	if userParam != "" {
		query = fmt.Sprintf("SELECT * FROM database.bookings WHERE UserID='%s'", userParam)
	} else if plotParam != "" {
		query = fmt.Sprintf("SELECT * FROM database.bookings WHERE PlotID='%s'", plotParam)
	} else {
		query = fmt.Sprintf("SELECT * FROM database.bookings")
	}

	bookings := getBookings(query)

	json.NewEncoder(w).Encode(bookings)
}

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
		err = results.Scan(&booking.BookingID, &booking.PlotID, &booking.UserID, &booking.StartDate, &booking.EndDate, &booking.LeaseCompleted)
		if err != nil {
			panic(err.Error())
		}

		tempBookings["bookings"] = append(tempBookings["bookings"], booking)
	}

	return tempBookings
}

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

			leaseCompleted, err := strconv.ParseBool(bookings["bookings"][0].LeaseCompleted)
			if err != nil {
				panic(err.Error())
			}

			// check if booking exists and has not yet been completed
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
				if newBooking.PlotID == "" || newBooking.UserID == "" || newBooking.StartDate == "" || newBooking.EndDate == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - All fields must be filled out and in JSON format"))
					return
				}

				// check if plot is available on desired dates
				if plotAvailable(newBooking.PlotID, newBooking.StartDate, newBooking.EndDate) {
					// establish connection to database
					db, err := sql.Open("mysql", connection)
					if err != nil {
						panic(err.Error())
					}
					defer db.Close()

					query := fmt.Sprintf("INSERT INTO bookings (PlotID, UserID, StartDate, EndDate, LeaseCompleted) VALUES ('%s', '%s', '%s', '%s', 'false')", newBooking.PlotID, newBooking.UserID, newBooking.StartDate, newBooking.EndDate)
					_, err = db.Query(query)
					if err != nil {
						panic(err.Error())
					}

					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Booking added on plot " + newBooking.PlotID))
				} else {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Desired dates are not available"))
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
				if editBooking.BookingID == "" || editBooking.PlotID == "" || editBooking.UserID == "" || editBooking.StartDate == "" || editBooking.EndDate == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - All fields must be filled out and in JSON format"))
					return
				}

				// check if plot is available for new dates
				if plotAvailable(editBooking.PlotID, editBooking.StartDate, editBooking.EndDate) {
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
					w.Write([]byte("422 - Desired dates are not available"))
				}

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply information in JSON format"))
			}
		}
	}
}

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

func plotAvailable(plotID string, startDate string, endDate string) (available bool) {
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
