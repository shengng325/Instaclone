package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	NotFoundError  = errors.New("model: resource not found")
	InvalidIdError = errors.New("models: ID provided was invalid")
)

//NewUserService Init
func NewUserService(connectionInfo string) (*UserService, error) {
	newDb, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	newDb.LogMode(true)
	us := &UserService{
		db: newDb,
	}
	return us, err
}

//UserService struct
type UserService struct {
	db *gorm.DB
}

//ByID will look up by the id provided
// 1 - user, nil
// 2 - nil, NotFoundError
// 3 - nil, OtherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

//ByEmail will look up by the id provided
// 1 - user, nil
// 2 - nil, NotFoundError
// 3 - nil, OtherError
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return NotFoundError
	}
	return err
}

//Create US
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

//Update US
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

//Delete US
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return InvalidIdError
	}
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	return us.db.Delete(&user).Error
}

//DestructiveReset Drop and auto migrate, only for dev
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

//Close for close at the end. Should call with defer
func (us *UserService) Close() error {
	return us.db.Close()
}

//User model
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
