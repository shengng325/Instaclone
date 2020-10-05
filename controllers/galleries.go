package controllers

import (
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
