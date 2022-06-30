package functions

import "net/http"

//Logout handles the logging out of a user session. It deletes the cookie in the session immediately
func Logout(res http.ResponseWriter, req *http.Request) {
	myCookie := &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
}

//LogoutAdmin handles the logging out of a Admin session. It deletes the cookie in the session immediately
func LogoutAdmin(res http.ResponseWriter, req *http.Request) {
	myCookieAdmin := &http.Cookie{
		Name:   "myCookieAdmin",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(res, myCookieAdmin)

	http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
}
