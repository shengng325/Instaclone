package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
	"lenslocked.com/models"
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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	services, err := models.NewServices(psqlInfo)
	must(err)
	defer services.Close()
	services.AutoMigrate()
	//services.DestructiveReset()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, r)
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/signup", usersC.SignupView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	r.Handle("/galleries",
		requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new",
		requireUserMw.Apply(galleriesC.NewGallery)).Methods("GET")
	r.HandleFunc("/galleries",
		requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}",
		galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")

	fmt.Println("Server running at :3000")
	http.ListenAndServe(":3000", userMw.Apply(r))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
