package functions

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Name     string `field:"Name"`
	Username string `field:"Username"`
	Email    string `field:"Email"`
}

type updateUsers struct {
	Name     string `field:"Name"`
	Username string `field:"Username"`
	Password []byte `field:"Password"`
	Email    string `field:"Email"`
}

var db *sql.DB
var tpl *template.Template

func init() {

	tpl = template.Must(template.ParseGlob("htmlTemplates/*"))

}
func connectUserDB() *sql.DB {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/goliveuserdb")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func SignUp(res http.ResponseWriter, req *http.Request) {

	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****SignUpHandler running*****")
	if req.Method == http.MethodGet {
		tpl.ExecuteTemplate(res, "signup.html", nil)
	} else if req.Method == http.MethodPost {

		name := req.FormValue("name")
		username := req.FormValue("username")
		password := req.FormValue("password")

		passwordVerification := false
		if 8 <= len(password) && len(password) < 60 {
			passwordVerification = true
		}

		if passwordVerification == false {
			tpl.ExecuteTemplate(res, "signup.html", "please check username and password criteria")
		}

		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 7)

		email := req.FormValue("email")
		query := fmt.Sprintf("INSERT INTO users (Name,Username, Password, Email ) VALUES ( '%v', '%v', '%v' ,'%v')", name, username, string(hashPassword), email)

		results, err := db.Query(query)
		defer results.Close()
		if err != nil {
			fmt.Println("Error in  Signup Query!")
		}

		http.Redirect(res, req, "/loginauth", 302)
	}

}

func LoginAuth(res http.ResponseWriter, req *http.Request) {
	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****loginAuthHandler running*****")
	if req.Method == http.MethodGet {
		tpl.ExecuteTemplate(res, "login.html", nil)
	} else {
		req.ParseForm()
		username := req.FormValue("username")
		password := req.FormValue("password")

		// retrieve password from db to compare (hash) with user supplied password's hash
		var hash string
		stmt := "SELECT Password FROM users WHERE Username = ?"
		row := db.QueryRow(stmt, username)
		err := row.Scan(&hash)

		if err != nil {
			fmt.Println("error selecting Hash in db by Username")
			tpl.ExecuteTemplate(res, "login.html", "check username and password")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		// returns nil on succcess
		if err == nil {
			fmt.Fprint(res, "You have successfully logged in :)")
			return
		}

		tpl.ExecuteTemplate(res, "login.html", "check username and password")
	}
}

func AllUsers(res http.ResponseWriter, req *http.Request) {
	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****AllUsesHandler running*****")
	if req.Method == http.MethodGet {
		results, err := db.Query("SELECT Name, Username, Email FROM goliveuserdb.users")
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

	stmt := fmt.Sprintf("DELETE FROM  goliveuserdb.users WHERE (`Username` = '%v')", username)
	result, err := db.Query(stmt)
	defer result.Close()
	if err != nil {
		panic(err)
	}

	tpl.ExecuteTemplate(res, "result.html", "User was Successfully Deleted")
}

func Update(res http.ResponseWriter, req *http.Request) {
	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****updateHandler running*****")
	req.ParseForm()
	username := req.FormValue("username")

	fmt.Println(username)
	query := fmt.Sprintf(`SELECT * FROM users WHERE Username = '%v'`, username)
	row := db.QueryRow(query)
	fmt.Println(row)
	var p updateUsers

	err := row.Scan(&p.Name, &p.Username, &p.Password, &p.Email)
	if err != nil {
		fmt.Println(err, "Scan error")
		http.Redirect(res, req, "/allusers", 307)
		return
	}
	tpl.ExecuteTemplate(res, "update.html", p)
}

func UpdateResult(res http.ResponseWriter, req *http.Request) {
	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****updateResultHandler running*****")
	req.ParseForm()
	name := req.FormValue("nameName")
	username := req.FormValue("userName")
	password := req.FormValue("passwordName")
	email := req.FormValue("emailName")
	stmt := fmt.Sprintf("UPDATE goliveuserdb.users SET name= '%v', Username = '%v', Password = '%v', Email = '%v' WHERE Username = '%v';", name, username, string(password), email, username)
	result, err := db.Query(stmt)
	result.Close()
	if err != nil {
		fmt.Println("error preparing stmt")
		panic(err)
	}
	tpl.ExecuteTemplate(res, "result.html", "User was successfully updated")
}
