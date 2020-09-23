package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"lenslocked.com/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null; unique_index"`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	//db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	user := models.User{
		Name:  "NTSSSSS",
		Email: "ntsss@mail.com",
	}
	err = us.Create(&user)
	if err != nil {
		panic(err)
	}

	//us.DestructiveReset()
	// user, err := us.ByID(1)
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println(user)

}
