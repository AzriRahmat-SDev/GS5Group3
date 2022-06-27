package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func getAllUsers(res http.ResponseWriter, req *http.Request) {
	uMap := makeUserMap("")
	json.NewEncoder(res).Encode(uMap)
}

func userHandler(res http.ResponseWriter, req *http.Request) {
	//the whole junk of if == "GET" goes here

	db := OpenUserDB()
	defer db.Close()

	params := mux.Vars(req)

	if req.Method == "GET" {
		if userExist(params["username"], "Username") {
			fmt.Println("Username: ", params["Username"])
			u := makeUserMap(params["username"])
			json.NewEncoder(res).Encode(u)
		} else {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("404 - No Username found from GET"))
		}
	}

	if req.Method == "DELETE" {
		if userExist(params["username"], "Username") {
			DeleteUsername(db, params["username"])
			res.WriteHeader(http.StatusAccepted)
			res.Write([]byte("202 - Username deleted: " + params["username"]))
		} else {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("404 - No such Username found"))
		}
	}

	if req.Header.Get("Content-type") == "application/json" {
		if req.Method == "POST" {
			var newUser Users
			reqBody, err := ioutil.ReadAll(req.Body)
			if err == nil {
				json.Unmarshal(reqBody, &newUser)
				if newUser.Username == "" {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply Username in JSON format"))
					return
				}
				if !userExist(params["username"], "Username") {
					InsertUser(db, newUser)
					fmt.Println("Insert was successful")
					res.WriteHeader(http.StatusCreated)
					res.Write([]byte("201 - Username added: " + params["username"]))
				} else {
					res.WriteHeader(http.StatusConflict)
					res.Write([]byte("409 - Duplicate Username"))
				}
			} else {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply Username information in JSON format"))
			}
		}
		if req.Method == "PUT" {
			var newUser Users
			reqBody, err := ioutil.ReadAll(req.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newUser)
				if newUser.Name == "" && newUser.Email == "" {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply user information in JSON format"))
					return
				}
				if !userExist(params["username"], "Username") {
					InsertUser(db, newUser)
					fmt.Println("Insert was successful")
					res.WriteHeader(http.StatusCreated)
					res.Write([]byte("201 - Information update: " + params["username"]))
				} else {
					if newUser.Name != "" {
						EditUserDisplayName(db, params["username"], newUser.Name)
					}
					if newUser.Email != "" {
						EditUserEmail(db, params["username"], newUser.Email)
					}
					if newUser.Username != "" {
						EditUsername(db, params["username"], newUser.Username)
					}
					res.Write([]byte("201-" + params["username"] + " has been updated" + "\nNew Display name: " + newUser.Name + "\nNew Username: " + newUser.Username + "\nNew Email: " + newUser.Email))
				}
			} else {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply Username information in JSON format"))
			}
		}
	}
}

func userExist(val string, column string) bool {
	db := OpenUserDB()
	defer db.Close()

	r := false

	s, err := db.Query("SELECT EXISTS (SELECT * FROM database.users WHERE " + column + "='" + val + "')")
	if err != nil {
		panic(err.Error())
	}
	for s.Next() {
		err = s.Scan(&r)
		if err != nil {
			panic(err.Error())
		}
	}
	return r
}

func makeUserMap(val string) UserMap {
	userMap := make(map[string]Users)
	db := OpenUserDB()
	defer db.Close()

	if val == "" {
		fmt.Println("Making Full map of users")
		query := fmt.Sprintf("SELECT Name,Username,Email FROM users")
		res, err := db.Query(query)
		if err != nil {
		}
		for res.Next() {
			var u Users
			res.Scan(&u.Name, &u.Username, &u.Email)
			userMap[u.Username] = u
		}
	} else {
		result, err := db.Query("SELECT Name,Username,Email from database.users WHERE Username = '" + val)
		if err != nil {
			fmt.Println(err)
		}
		for result.Next() {
			var u Users
			err := result.Scan(&u.Name, &u.Username, &u.Email)
			if err != nil {
				fmt.Println(err.Error())
			}
			userMap["Username"] = u
		}
	}
	return userMap
}
