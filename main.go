package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1> Home. Welcome! </h1>")
}

func ContactsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1> Contact page </h1>")
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
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/contacts", ContactsHandler)
	router.GET("/faq", FaqHandler)
	router.NotFound = http.HandlerFunc(NotFoundHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}
