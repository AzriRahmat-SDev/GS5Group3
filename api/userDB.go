package api

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Users struct {
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Password []byte `json:"Password"`
	Email    string `json:"Email"`
}

type UserMap map[string]Users

var userList []Users

func OpenUserDB() *sql.DB {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connected to user db")
	return db
}

func InsertUser(db *sql.DB, u Users) {
	query := fmt.Sprintf("INSERT INTO users (Name, Username,Password, Email) VALUES ('%s', '%s','%s',  '%s')", u.Name, u.Username, u.Password, u.Email)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nSuccessful insert user @ '%s'", u)
	}
}

func EditUserDisplayName(db *sql.DB, Username string, Name string) {
	query := fmt.Sprintf("Update users SET Name = '%s' WHERE Username = '%s'", Name, Username)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nSuccessful update Users Display name @ '%s' with '%s'", Username, Name)
	}
}

func EditUserEmail(db *sql.DB, Username string, Email string) {
	query := fmt.Sprintf("Update users SET Email = '%s' WHERE Username = '%s'", Email, Username)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nSuccessful update Users Email @ '%s' with '%s'", Username, Email)

	}
}

func EditUsername(db *sql.DB, Username string, newUsername string) {
	query := fmt.Sprintf("Update users SET Username = '%s' WHERE Username = '%s'", newUsername, Username)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nSuccessful updated User's old Username = '%s' with new Username '%s'", Username, newUsername)
	}
}

func DeleteUsername(db *sql.DB, Username string) {
	stmt := fmt.Sprintf("DELETE FROM users WHERE Username='%s'", Username)
	_, err := db.Query(stmt)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nSuccessfully deleted username @ '%s'", Username)
	}
}
