package functions

import (
	"net/http"
)

func alreadyLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	_, err := req.Cookie("myCookie")

	if err != nil {

		http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
		return false
	}

	return true
}
