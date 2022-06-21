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
	var err error
	tpl = template.Must(template.ParseGlob("htmlTemplates/*"))
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/goliveuserdb")
	if err != nil {
		panic(err.Error())

	}
}

func SignUp(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Check signup")
	if req.Method == http.MethodGet {
		fmt.Println("Check get")
		tpl.ExecuteTemplate(res, "signup.html", nil)
	} else if req.Method == http.MethodPost {
		fmt.Println("Check post")

		name := req.FormValue("name")
		username := req.FormValue("username")
		password := req.FormValue("password")
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 7)

		email := req.FormValue("email")
		query := fmt.Sprintf("INSERT INTO users (Name,Username, Password, Email ) VALUES ( '%v', '%v', '%v' ,'%v')", name, username, string(hashPassword), email)

		_, err := db.Query(query)

		if err != nil {
			panic(err.Error())
		}

		http.Redirect(res, req, "/loginauth", 302)
	}

}

func LoginAuth(res http.ResponseWriter, req *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	if req.Method == http.MethodGet {
		fmt.Println("Check get")
		tpl.ExecuteTemplate(res, "login.html", nil)
	} else {
		req.ParseForm()
		username := req.FormValue("username")
		password := req.FormValue("password")
		fmt.Println("username:", username, "password:", password)
		// retrieve password from db to compare (hash) with user supplied password's hash
		var hash string
		stmt := "SELECT Password FROM users WHERE Username = ?"
		row := db.QueryRow(stmt, username)
		err := row.Scan(&hash)
		fmt.Println("hash from db:", hash)
		if err != nil {
			fmt.Println("error selecting Hash in db by Username")
			tpl.ExecuteTemplate(res, "login.html", "check username and password")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		// returns nill on succcess
		if err == nil {
			fmt.Fprint(res, "You have successfully logged in :)")
			return
		}
		fmt.Println("incorrect password")
		tpl.ExecuteTemplate(res, "login.html", "check username and password")
	}
}

func AllUsers(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		results, err := db.Query("SELECT Name, Username, Email FROM goliveuserdb.users")
		if err != nil {
			panic("Error in results")
		}
		fmt.Println("Successfully Selected all")
		defer results.Close()
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
	fmt.Println("*****deleteHandler running*****")
	req.ParseForm()
	id := req.FormValue("username")
	del, err := db.Prepare("DELETE FROM  goliveuserdb.users WHERE (`Username` = ?);")
	if err != nil {
		panic(err)
	}
	defer del.Close()
	var result sql.Result
	result, err = del.Exec(id)
	rowsAff, _ := result.RowsAffected()
	fmt.Println("rowsAff:", rowsAff)

	if err != nil || rowsAff != 1 {
		fmt.Fprint(res, "Error deleting product")
		return
	}

	fmt.Println("err:", err)
	tpl.ExecuteTemplate(res, "result.html", "User was Successfully Deleted")
}

func Update(res http.ResponseWriter, req *http.Request) {
	fmt.Println("*****updateHandler running*****")
	req.ParseForm()
	username := req.FormValue("username")
	// query := fmt.Sprintf(`SELECT COUNT(*) FROM Users WHERE Id = '%s'`, userID)
	// row := DB.QueryRow(query)
	// if err != nil {
	// 	return err.Error()
	// }
	fmt.Println(username)
	query := fmt.Sprintf(`SELECT * FROM users WHERE Username = '%v'`, username)
	row := db.QueryRow(query)
	fmt.Println(row)
	var p updateUsers
	// func (r *Row) Scan(dest ...interface{}) error
	err := row.Scan(&p.Name, &p.Username, &p.Password, &p.Email)
	if err != nil {
		fmt.Println(err, "Scan error")
		http.Redirect(res, req, "/browse", 307)
		return
	}
	tpl.ExecuteTemplate(res, "update.html", p)
}

func UpdateResult(res http.ResponseWriter, req *http.Request) {
	fmt.Println("*****updateResultHandler running*****")
	req.ParseForm()
	name := req.FormValue("nameName")
	username := req.FormValue("userName")
	password := req.FormValue("passwordName")
	email := req.FormValue("emailName")
	upStmt := fmt.Sprintf("UPDATE goliveuserdb.users SET name= '%v', Username = '%v', Password = '%v', Email = '%v' WHERE Username = '%v';", name, username, string(password), email, username)
	fmt.Printf("UPDATE goliveuserdb.users SET name= '%v', Username = '%v', Password = '%v', Email = '%v' WHERE Username = '%v';", name, username, password, email, username)

	// func (db *DB) Prepare(query string) (*Stmt, error)
	_, err := db.Query(upStmt)
	if err != nil {
		fmt.Println("error preparing stmt")
		panic(err)
	}
	fmt.Println("db.Prepare err:", err)
	tpl.ExecuteTemplate(res, "result.html", "User was successfully updated")
}
func AccountManagement() {

}
