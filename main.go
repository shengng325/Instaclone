package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"lenslocked.com/views"
)

var homeView *views.View
var contactView *views.View

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Templates.ExecuteTemplate(w, homeView.Layout, nil)
	if err != nil {
		panic(err)
	}
}

func ContactsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Templates.ExecuteTemplate(w, contactView.Layout, nil)
	if err != nil {
		panic(err)
	}

}

func FaqHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1> FAQ </h1>")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1> 404 NOT FOUND!!! </h1>")
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/contacts", ContactsHandler)
	router.GET("/faq", FaqHandler)
	router.NotFound = http.HandlerFunc(NotFoundHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}
