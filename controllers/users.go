package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/views"
)

//init user
func InitUser() *User {
	return &User{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}

type User struct {
	NewView *views.View
}

func (u *User) Handler(w http.ResponseWriter, r *http.Request) {
	err := u.NewView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	err := parseForm(r, &form)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)

}
