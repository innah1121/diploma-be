package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

func Connect() (*sql.DB, error) {
	// Open up our database connection.
	db, err := sql.Open("mysql", "root:password1@tcp(127.0.0.1:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	// defer the close after has finished executing
	defer db.Close()

	return db, nil
}

type User struct {
    Id    int
    Username  string
    Password string
}

func Insert(username string, password string) {

    insForm, err := db.Prepare("INSERT INTO User(username, password) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(username, password)
        fmt.Println("INSERT: Username: " + username + " | Password: " + password)

}

func GetByUsername(username string) {

    selDB, err := db.Query("SELECT * FROM Employee WHERE username=?", username)
    if err != nil {
        panic(err.Error())
    }
	fmt.Println(selDB)

}

func GetByUsernameAndPassword(username string,password string) {

    selDB, err := db.Query("SELECT * FROM Employee WHERE username=? AND password=?", username, password)
    if err != nil {
        panic(err.Error())
    }
	fmt.Println(selDB)

}

