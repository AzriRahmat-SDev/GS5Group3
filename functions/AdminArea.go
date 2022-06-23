package functions

import (
	"fmt"
	"net/http"
)

type Users struct {
	Name     string `field:"Name"`
	Username string `field:"Username"`
	Email    string `field:"Email"`
}

func AllUsers(res http.ResponseWriter, req *http.Request) {
	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****AllUsesHandler running*****")
	if req.Method == http.MethodGet {
		results, err := db.Query("SELECT Name, Username, Email FROM users")
		defer results.Close()
		if err != nil {
			panic("Error in Allusers Query")
		}

		var userArr []Users

		for results.Next() {
			var user Users
			err := results.Scan(&user.Name, &user.Username, &user.Email)
			if err != nil {
				panic("Error in scan")
			}

			userArr = append(userArr, user)
		}

		tpl.ExecuteTemplate(res, "allusers.html", userArr)

	}

}

func DeleteRecord(res http.ResponseWriter, req *http.Request) {
	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****deleteHandler running*****")
	req.ParseForm()
	username := req.FormValue("username")

	stmt := fmt.Sprintf("DELETE FROM users WHERE (`Username` = '%v')", username)
	result, err := db.Query(stmt)
	defer result.Close()
	if err != nil {
		panic(err)
	}

	tpl.ExecuteTemplate(res, "result.html", "User was Successfully Deleted")
}
