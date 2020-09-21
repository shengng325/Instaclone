package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

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

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	var form SignupForm
	dec := schema.NewDecoder()
	dec.Decode(&form, r.PostForm)
	fmt.Fprintln(w, form)

}
