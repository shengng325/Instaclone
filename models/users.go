package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
	"lenslocked.com/rand"

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
const hmacSecretKey = "secret-hmac-key"

//User model
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

//UserDB is
type UserDB interface {
	//Methods got querying users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//Altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	Close() error

	AutoMigrate() error
	DestructiveReset() error
}

// UserService defines all of the methods we need to
// interact with the User resource.
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

//NewUserService Init
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	// newDb, err := gorm.Open("postgres", connectionInfo)
	// if err != nil {
	// 	return nil, err
	// }
	// newDb.LogMode(true)
	// hmac := hash.NewHMAC(hmacSecretKey)
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	us := &userService{
		UserDB: uv,
	}
	return us, err
}

var _ UserService = &userService{}

//UserService struct
type userService struct {
	UserDB
}

//Authenticate use to authenticate email and pw
func (us *userService) Authenticate(email, password string) (*User, error) {
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

type validatorFunc func(*User) error

func runValidatorFunc(user *User, fns ...validatorFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

//userValidator struct
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runValidatorFunc(&user,
		uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

//Create US
func (uv *userValidator) Create(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	err := runValidatorFunc(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runValidatorFunc(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

//Delete US
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return InvalidIdError
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	ug := &userGorm{
		db: db,
	}
	return ug, err
}

//UserService struct
type userGorm struct {
	db *gorm.DB
}

//ByID will look up by the id provided
// 1 - user, nil
// 2 - nil, NotFoundError
// 3 - nil, OtherError
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//ByEmail will look up by the id provided
// 1 - user, nil
// 2 - nil, NotFoundError
// 3 - nil, OtherError
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	db := ug.db.Where("remember_token = ?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Create US
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

//Update US
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//Delete US
func (ug *userGorm) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	return ug.db.Delete(&user).Error
}

//DestructiveReset Drop and auto migrate, only for dev
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

//AutoMigrate will auto migrate the users table
func (ug *userGorm) AutoMigrate() error {
	err := ug.db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

//Close for close at the end. Should call with defer
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return NotFoundError
	}
	return err
}
