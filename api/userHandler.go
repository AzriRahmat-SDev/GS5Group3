package api

import (
	"encoding/json"
	"net/http"
)

func GetAllUsers(res http.ResponseWriter, req *http.Request) {
	db := OpenUserDB()
	defer db.Close()
	populateUserData(db)
	json.NewEncoder(res).Encode(userMap)
}

//settle user handler and take reference from victors code
