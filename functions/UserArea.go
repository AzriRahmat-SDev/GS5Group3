package functions

import (
	"fmt"
	"net/http"
)

func UserArea(res http.ResponseWriter, req *http.Request) {
	// pull user info
	cookie, err := req.Cookie("myCookie")
	if err != nil {
		http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
		return
	}
	user := getUser(cookie)

	// pull bookings for plot
	allBookings := callBookingsAPI("byUser", user.Username)

	// separate current and expired bookings
	var currentBookings, expiredBookings bookings

	for _, v := range allBookings.Bookings {
		if v.LeaseCompleted == "true" {
			expiredBookings.Bookings = append(expiredBookings.Bookings, v)
		} else if v.LeaseCompleted == "false" {
			currentBookings.Bookings = append(currentBookings.Bookings, v)
		} else {
			fmt.Printf("LeaseCompleted field does not exist for this entry: %v", v)
		}
	}

	allInfo := map[string]allInfo{
		"allInfo": {
			Username:      user.Username,
			Name:          user.Name,
			Email:         user.Email,
			CurrentLeases: currentBookings,
			ExpiredLeases: expiredBookings,
		},
	}

	tpl.ExecuteTemplate(res, "userarea.gohtml", allInfo)
}
