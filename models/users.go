package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	NotFoundError = errors.New("model: resource not found")
)

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

type UserService struct {
	db *gorm.DB
}

//ByID will look up by the id provided
// 1 - user, nil
// 2 - nil, NotFoundError
// 3 - nil, OtherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, NotFoundError
	default:
		return nil, err
	}
}

func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

func (us *UserService) Close() error {
	return us.db.Close()
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
