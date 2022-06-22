package main

import (
	"GitHub/GS5Group3/functions"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("htmlTemplates/*"))
}

func main() {
	r := mux.NewRouter()
	//viewing page. Has the option to go into adminlogin
	//r.HandleFunc("/", functions.IndexPage)

	//account funtions

	r.HandleFunc("/signup", functions.SignUp)

	r.HandleFunc("/loginauth", functions.LoginAuth)
	r.HandleFunc("/allusers", functions.AllUsers)
	r.HandleFunc("/delete/", functions.DeleteRecord)
	r.HandleFunc("/update/", functions.Update)
	r.HandleFunc("/updateresult/", functions.UpdateResult)
	// r.HandleFunc("/accountmanagement", functions.AccountManagement)
	r.HandleFunc("/newbooking/", functions.NewBooking)
	r.HandleFunc("/editbooking/", functions.EditBooking)
	r.HandleFunc("/deletebooking/", functions.DeleteBooking)

	// //User Area functions
	// r.HandleFunc("/bookspace", functions.BookSpace)

	// //Admin Area
	// r.HandleFunc("/adminlogin", functions.AdminLogin)
	// r.HandleFunc("/adminarea", functions.AdminArea)

	log.Fatal(http.ListenAndServe(":8080", r))
}
