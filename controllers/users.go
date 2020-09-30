package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

//init user
func InitUser(us models.UserService) *User {
	return &User{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type User struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
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
	u.signIn(w, &user)
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

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.NotFoundError:
			fmt.Fprintln(w, "Invalid email address")
		case models.ErrIncorrectPassword:
			fmt.Fprintln(w, "Invalid password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	u.signIn(w, user)
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
