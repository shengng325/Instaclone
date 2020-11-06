package controllers

import (
	"net/http"
	"strconv"

	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

func NewStatic(gs models.GalleryService) *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
		gs:      gs,
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
	gs      models.GalleryService
}

func (s *Static) HomeRedirect(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	// galleryId := galleries[0]
	if user != nil {
		galleries, _ := s.gs.ByUserID(user.ID)
		if len(galleries) > 0 {
			galleryId := galleries[0].ID
			http.Redirect(w, r, "/galleries/"+strconv.Itoa(int(galleryId)), http.StatusFound)
		} else {
			http.Redirect(w, r, "/galleries", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
