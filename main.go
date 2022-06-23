package main

import (
	"GS5Group3/api"
	"GS5Group3/functions"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("htmlTemplates/*"))
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

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
	r.HandleFunc("/homepage/", functions.Homepage)
	// r.HandleFunc("/accountmanagement", functions.AccountManagement)
	// r.HandleFunc("/bookings", functions.Bookings) // Holds both current and booking history
	r.HandleFunc("/newbooking/", functions.NewBooking)
	r.HandleFunc("/editbooking/", functions.EditBooking)
	r.HandleFunc("/deletebooking/", functions.DeleteBooking)
	r.HandleFunc("/logout", functions.Logout)

	// //User Area functions
	// r.HandleFunc("/bookspace", functions.BookSpace)

	// //Admin Area
	// r.HandleFunc("/adminlogin", functions.AdminLogin)
	// r.HandleFunc("/adminarea", functions.AdminArea)

	go func() {
		fmt.Println("Starting API for venue and booking...")
		api.StartServer()
		wg.Done()
	}()

	go func() {
		fmt.Println("Starting server for users...")
		log.Fatal(http.ListenAndServe(":8080", r))
		wg.Done()
	}()

	wg.Wait()
}
