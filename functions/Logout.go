package functions

import "net/http"

func logout(res http.ResponseWriter, req *http.Request) {
	myCookie := &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(res, myCookie)
	http.Redirect(res, req, "/homepage/", http.StatusSeeOther)
}
