package functions

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func connectUserDB() *sql.DB {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(localhost:32769)/database")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("connected to user db")
	return db
}
