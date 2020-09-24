package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

//init user
func InitUser(us *models.UserService) *User {
	return &User{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type User struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

func (u *User) Handler(w http.ResponseWriter, r *http.Request) {
	err := u.NewView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	err := parseForm(r, &form)
	if err != nil {
		panic(err)
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	err = u.us.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintln(w, user)

}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//Login verify email and pw
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	err := parseForm(r, &form)
	if err != nil {
		panic(err)
	}
}
