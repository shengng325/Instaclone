package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// var id int
	// var name, email string
	// row := db.QueryRow(`
	// 	INSERT INTO users(name, email)
	// 	VALUES($1, $2)
	// 	RETURNING id`,
	// 	"Jons Calhoun", "jons@calhoun.io")
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("return id...", id)

	// row := db.QueryRow(`
	// 	SELECT id, name, email
	// 	FROM users
	// 	WHERE id=$1`, 1)
	// err = row.Scan(&id, &name, &email)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("return id...", id, "...name...", name, "...email...", email)

	// type User struct {
	// 	id    int
	// 	name  string
	// 	email string
	// }
	// var users []User
	// rows, err := db.Query(`
	// 	SELECT id, name, email
	// 	FROM users`)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var user User
	// 	err = rows.Scan(&user.id, &user.name, &user.email)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	users = append(users, user)
	// }
	// if rows.Err() != nil {
	// 	panic(err)
	// }
	// fmt.Println(users)

	var id int
	for i := 1; i < 6; i++ {
		// Create some fake data
		userId := 1
		if i > 3 {
			userId = 2
		}
		amount := 1000 * i
		description := fmt.Sprintf("USB-C Adapter x%d", i)
		_, err = db.Exec(`
			INSERT INTO orders (user_id, amount, description)
			VALUES($1, $2, $3)`,
			userId, amount, description)
		if err != nil {
			panic(err)
		}
		fmt.Println("Created an order with the ID:", id)
	}

	db.Close()
}
