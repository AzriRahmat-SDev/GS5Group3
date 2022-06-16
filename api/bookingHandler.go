package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func getBookings(w http.ResponseWriter, r *http.Request) {
	// establish connection to database
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM database.bookings")
	if err != nil {
		panic(err.Error())
	}

	// create variable to store all friends
	bookings := map[string][]booking{
		"bookings": {},
	}

	for results.Next() {
		var booking booking
		err = results.Scan(&booking.BookingID, &booking.PlotID, &booking.UserID, &booking.StartDate, &booking.EndDate, &booking.LeaseCompleted)
		if err != nil {
			panic(err.Error())
		}

		bookings["bookings"] = append(bookings["bookings"], booking)
	}

	json.NewEncoder(w).Encode(bookings)
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingParam := params["BookingID"]

	// establish connection to database
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if r.Method == "GET" {
		if bookingExists(db, bookingParam) {
			results, err := db.Query("SELECT * FROM database.bookings WHERE BookingID = '" + bookingParam + "' LIMIT 1")
			if err != nil {
				panic(err.Error())
			}

			// create variable to store booking
			bookings := map[string][]booking{
				"bookings": {},
			}

			for results.Next() {
				var booking booking
				err = results.Scan(&booking.BookingID, &booking.PlotID, &booking.UserID, &booking.StartDate, &booking.EndDate, &booking.LeaseCompleted)
				if err != nil {
					panic(err.Error())
				}

				bookings["bookings"] = append(bookings["bookings"], booking)
			}

			json.NewEncoder(w).Encode(bookings)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Booking not found"))
		}
	}

	if r.Method == "DELETE" {
		results, err := db.Query("SELECT * FROM database.bookings WHERE BookingID = '" + bookingParam + "' LIMIT 1")
		if err != nil {
			panic(err.Error())
		}

		// create variable to store booking
		bookings := map[string][]booking{
			"bookings": {},
		}

		for results.Next() {
			var booking booking
			err = results.Scan(&booking.BookingID, &booking.PlotID, &booking.UserID, &booking.StartDate, &booking.EndDate, &booking.LeaseCompleted)
			if err != nil {
				panic(err.Error())
			}

			bookings["bookings"] = append(bookings["bookings"], booking)
		}

		leaseCompleted, err := strconv.ParseBool(bookings["bookings"][0].LeaseCompleted)
		if err != nil {
			panic(err.Error())
		}

		// check if booking exists and has not yet been completed
		if bookingExists(db, bookingParam) && !leaseCompleted {
			query := fmt.Sprintf("DELETE FROM database.bookings WHERE BookingID='%s'", bookingParam)

			_, err := db.Query(query)
			if err != nil {
				panic(err.Error())
			}

			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Booking canceled: " + bookingParam))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Booking not found or has already been completed"))
		}
	}
}

func bookingExists(db *sql.DB, booking string) (exists bool) {
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
