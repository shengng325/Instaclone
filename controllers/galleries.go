package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

//init user
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		NewGallery: views.NewView("bootstrap", "galleries/new"),
		gs:         gs,
	}
}

type Galleries struct {
	NewGallery *views.View
	gs         models.GalleryService
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.NewGallery.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Println("User: ", user)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewGallery.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
