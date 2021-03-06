package functions

import (
	"fmt"
	"net/http"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type updateUsers struct {
	Name     string `field:"Name"`
	Username string `field:"Username"`
	Password []byte `field:"Password"`
	Email    string `field:"Email"`
}

/*SignUp handles either a GET or POST request.
GET request returns a blank form
PUT request will request a display name,username, password and email values which will be then stored in the user database.
*/
func SignUp(res http.ResponseWriter, req *http.Request) {

	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****SignUpHandler running*****")
	if req.Method == http.MethodGet {
		tpl.ExecuteTemplate(res, "signup.html", nil)
	} else if req.Method == http.MethodPost {

		name := req.FormValue("name")
		username := req.FormValue("username") //have to add check for unique username
		password := req.FormValue("password")

		//password verification
		cfmpassword := req.FormValue("cfmpassword")
		if cfmpassword != password {
			tpl.ExecuteTemplate(res, "signup.html", "password did not match!")
			return
		}

		var pswdLowercase, pswdUppercase, passwordVerification, pswdNumber, pswdSpecial bool
		for _, char := range password {
			if unicode.IsLower(char) {
				pswdLowercase = true
			} else if unicode.IsUpper(char) {
				pswdUppercase = true
			} else if unicode.IsNumber(char) {
				pswdNumber = true
			} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
				pswdSpecial = true
			}
		}
		if 8 <= len(password) {
			passwordVerification = true
		}
		if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !passwordVerification {
			tpl.ExecuteTemplate(res, "signup.html", "Please check Username and Password criteria")
			return
		}
		//end password verification

		//hashing password for more security incase hackers get our userslist
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 7)

		email := req.FormValue("email")
		query := fmt.Sprintf("INSERT INTO users (Name,Username, Password, Email ) VALUES ( '%v', '%v', '%v' ,'%v')", name, username, string(hashPassword), email)

		results, err := db.Query(query)
		defer results.Close()
		if err != nil {
			fmt.Println("Error in  Signup Query!")
		}
		InfoLogger.Printf("User %s created.", username)
		http.Redirect(res, req, "/loginauth", 302)
	}

}

/*LoginAuth handles the authentication of the user details in the user database and logs the user in.
Upon successful login,  a user session is created using a cookie.
*/
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

		//Admin Area
		if username == "Admin" {
			if password == "Admin123!@#" {
				myCookieAdmin := &http.Cookie{
					Name:   "myCookieAdmin",
					Value:  username,
					MaxAge: 3600,
				}
				http.SetCookie(res, myCookieAdmin)
				http.Redirect(res, req, "/allusers", 303)
				return
			}
		}
		// retrieve password from db to compare (hash) with user supplied password's hash
		var hash string
		stmt := "SELECT Password FROM users WHERE Username = ?"
		row := db.QueryRow(stmt, username)
		err := row.Scan(&hash)

		if err != nil {
			fmt.Println("error selecting Hash in db by Username")
			tpl.ExecuteTemplate(res, "login.html", "Either Username or Password was incorrect")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		// returns nil on succcess
		if err == nil {
			myCookie := &http.Cookie{
				Name:   "myCookie",
				Value:  username,
				MaxAge: 3600,
			}
			http.SetCookie(res, myCookie)
			http.Redirect(res, req, "/venues/", 303)
			return
		}

		tpl.ExecuteTemplate(res, "login.html", "Check Username and Password")
	}
}

//Update handles the scanning of the username in the database
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

//UpdateResults handles the updating of the users values
func UpdateResult(res http.ResponseWriter, req *http.Request) {
	db := connectUserDB()
	defer db.Close()
	fmt.Println("*****updateResultHandler running*****")
	req.ParseForm()
	name := req.FormValue("nameName")
	username := req.FormValue("userName")
	password := req.FormValue("passwordName")
	email := req.FormValue("emailName")
	stmt := fmt.Sprintf("UPDATE users SET name= '%v', Username = '%v', Password = '%v', Email = '%v' WHERE Username = '%v';", name, username, string(password), email, username)
	result, err := db.Query(stmt)
	result.Close()
	if err != nil {
		fmt.Println("error preparing stmt")
		panic(err)
	}
	tpl.ExecuteTemplate(res, "result.html", "User was successfully updated")
}
