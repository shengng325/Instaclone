package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServicesConfig func(*Services) error

// WithGorm will open a GORM connection with the provided
// info and attach it to the Services type if there aren't
// any errors.
func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

// WithLogMode will set the LogMode of the current GORM DB
// object associated with the Services type. It is assumed
// that the DB object will already exist and be initialized
// properly, so you will want to call something like
// WithGorm before this config function.
func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

// WithUser will use the existing GORM DB connection of the
// Services object along with the provided pepper and hmacKey
// to build and set a UserService.
func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

// WithGallery will use the existing GORM DB connection of
// the Services object to build and set a GalleryService.
func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

// WithImage will build and set an ImageService.
func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

// NewServices now will accept a list of config functions to
// run. Each function will accept a pointer to the current
// Services object as its only argument and will edit that
// object inline and return an error if there is one. Once
// we have run all configs we will return the Services object.
func NewServices(cfgs ...ServicesConfig) (*Services, error) {

	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil

}

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

//Close for close at the end. Should call with defer
func (s *Services) Close() error {
	return s.db.Close()
}

//DestructiveReset Drop and auto migrate, only for dev
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

//AutoMigrate will auto migrate the users table
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
