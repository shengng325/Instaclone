package models

import "strings"

const (
	NotFoundError        modelError = "models: resource not found"
	InvalidIdError       modelError = "models: ID provided was invalid"
	ErrIncorrectPassword modelError = "models: incorrect password provided"
	ErrEmailRequired     modelError = "models: email address is required"
	ErrEmailInvalid      modelError = "models: email address is invalid"
	ErrEmailTaken        modelError = "models: Email is already taken"
	ErrPasswordTooShort  modelError = "models: password must be at least 8 characters long"
	ErrPasswordRequired  modelError = "models: password is required"
	ErrRememberTooShort  modelError = "models: remember token must be at least 32 bytes"
	ErrRememberRequired  modelError = "models: remember is required"
)

type PublicError interface {
	error
	Public() string
}

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
