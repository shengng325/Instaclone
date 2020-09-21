package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/views"
)

//init user
func InitUser() *User {
	return &User{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
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

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Temp response")
}
