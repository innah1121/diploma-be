package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	// Open up our database connection.
	db, err := sql.Open("mysql", "root:password1@tcp(127.0.0.1:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	// defer the close after has finished executing
	// defer db.Close()

	return db, nil
}

type User struct {
    Id    int
    Username  string
    Password string
}

func Insert(username string, password string) {
    db , _ := Connect()
    insForm, err := db.Prepare("INSERT INTO User(username, password) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(username, password)
        log.Println("INSERT: Username: " + username + " | Password: " + password)
	defer db.Close()
}

