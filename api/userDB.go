package api

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Users struct {
	Name     string `field:"Name"`
	Username string `field:"Username"`
	Email    string `field:"Email"`
}

var db *sql.DB

var userMap map[string]Users
var userList []Users

func OpenUserDB() *sql.DB {
	var err error
	db, err = sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connected to user db")
	return db
}
func DeleteRecord(res http.ResponseWriter, req *http.Request) {

	db := OpenUserDB()
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

}

func populateUserData(db *sql.DB) {
	for k := range userMap {
		delete(userMap, k)
	}

	results, err := db.Query("Select Name, Username, Email FROM users")
	if err != nil {
		fmt.Println(err.Error())
	}
	for results.Next() {
		var u Users
		err := results.Scan(&u.Name, &u.Username, &u.Email)
		if err != nil {
			fmt.Println(err.Error())
		}
		userList = append(userList, u)
	}
}
