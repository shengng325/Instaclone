package views

import "lenslocked.com/models"

const (
	LevelError   = "danger"
	LevelWarning = "warning"
	LevelInfo    = "info"
	LevelSuccess = "success"
)

var (
	AlertGeneric = ("Something went " +
		"wrong. Please try again, and contact us if the " +
		"problem persists.")
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(models.PublicError); ok {
		d.Alert = &Alert{
			Level:   LevelError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   LevelError,
			Message: AlertGeneric,
		}
	}
}
