package functions

import "net/http"

func Logout(res http.ResponseWriter, req *http.Request) {
	myCookie := &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
}

func LogoutAdmin(res http.ResponseWriter, req *http.Request) {
	myCookieAdmin := &http.Cookie{
		Name:   "myCookieAdmin",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(res, myCookieAdmin)

	http.Redirect(res, req, "/loginauth", http.StatusSeeOther)
}
