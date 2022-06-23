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
