package controllers

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

//init user
func InitUser(us models.UserService) *User {
	return &User{
		SignupView: views.NewView("bootstrap", "users/new"),
		LoginView:  views.NewView("bootstrap", "users/login"),
		us:         us,
	}
}

type User struct {
	SignupView *views.View
	LoginView  *views.View
	us         models.UserService
}

func (u *User) Handler(w http.ResponseWriter, r *http.Request) {
	u.SignupView.Render(w, nil)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	err := parseForm(r, &form)
	if err != nil {
		vd.SetAlert(err)
		u.SignupView.Render(w, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	err = u.us.Create(&user)
	if err != nil {
		vd.SetAlert(err)
		u.SignupView.Render(w, vd)
	}
	err = u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
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
	vd := views.Data{}
	err := parseForm(r, &form)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.NotFoundError:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, user)
}

func (u *User) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

//CookieTest display cookies set
func (u *User) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
