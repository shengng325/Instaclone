package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"lenslocked.com/controllers"
)

// var homeView *views.View
// var contactView *views.View

// func HomeHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	must(homeView.Render(w, nil))

// }

// func ContactsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	must(contactView.Render(w, nil))
// }

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1> 404 NOT FOUND!!! </h1>")
}

func main() {
	// homeView = views.NewView("bootstrap", "views/home.gohtml")
	// contactView = views.NewView("bootstrap", "views/contact.gohtml")
	staticC := controllers.NewStatic()
	usersC := controllers.InitUser()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	//r.HandleFunc("/", HomeHandler).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	//r.HandleFunc("/contact", ContactsHandler).Methods("GET")
	r.HandleFunc("/signup", usersC.Handler).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
