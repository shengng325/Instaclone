package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserID uint   `gorm: "not_null; index"`
	Title  string `gorm: "not_null"`
}

type GalleryService interface{}

type GalleryDB interface{}
