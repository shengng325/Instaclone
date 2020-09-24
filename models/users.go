package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	NotFoundError        = errors.New("model: resource not found")
	InvalidIdError       = errors.New("models: ID provided was invalid")
	InvalidPasswordError = errors.New("models: incorrect password provided")
	//InvalidEmailError = errors.New("models: incorrect email address provided")
)

const userPwPepper = "secret-random-string"

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

//Authenticate use to authenticate email and pw
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, InvalidPasswordError
		default:
			return nil, err
		}
	}
	return foundUser, nil
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
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
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
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

//AutoMigrate will auto migrate the users table
func (us *UserService) AutoMigrate() error {
	err := us.db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

//Close for close at the end. Should call with defer
func (us *UserService) Close() error {
	return us.db.Close()
}

//User model
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}
