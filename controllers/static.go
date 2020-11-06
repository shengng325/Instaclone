package controllers

import (
	"net/http"

	"lenslocked.com/context"
	"lenslocked.com/views"
)

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}

func (s *Static) HomeRedirect(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user != nil {
		http.Redirect(w, r, "/galleries/1", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
