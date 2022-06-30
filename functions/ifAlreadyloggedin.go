package functions

import (
	"net/http"
)

//alreadyLoggedIn checks for the cookie in the user session and allowing the user to have access to the appropriate handler
func alreadyLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	_, err := req.Cookie("myCookie")

	if err != nil {
		tpl.ExecuteTemplate(res, "restricted.html", "You don't belong here")
		// http.Redirect(res, req, "/restricted", http.StatusSeeOther)
		return false
	}

	return true
}

//alreadyLoggedInAdmin checks for the cookie in the Admin session and allowing the Admin to have access to the appropriate handler
func alreadyLoggedInAdmin(res http.ResponseWriter, req *http.Request) bool {
	_, err := req.Cookie("myCookieAdmin")

	if err != nil {
		tpl.ExecuteTemplate(res, "restricted.html", "You don't belong here")
		// http.Redirect(res, req, "/restricted", http.StatusSeeOther)
		return false
	}

	return true
}
