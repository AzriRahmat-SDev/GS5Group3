package functions

import (
	"net/http"
)

func alreadyLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	_, err := req.Cookie("myCookie")

	if err != nil {
		tpl.ExecuteTemplate(res, "restricted.html", "You don't belong here")
		// http.Redirect(res, req, "/restricted", http.StatusSeeOther)
		return false
	}

	return true
}

func alreadyLoggedInAdmin(res http.ResponseWriter, req *http.Request) bool {
	_, err := req.Cookie("myCookieAdmin")

	if err != nil {
		tpl.ExecuteTemplate(res, "restricted.html", "You don't belong here")
		// http.Redirect(res, req, "/restricted", http.StatusSeeOther)
		return false
	}

	return true
}
