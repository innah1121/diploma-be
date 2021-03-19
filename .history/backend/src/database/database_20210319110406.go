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

type Employee struct {
    Id    int
    Username  string
    Password string
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db, _ := Connect()
    if r.Method == "POST" {
        name := r.FormValue("name")
        city := r.FormValue("city")
        insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, city)
        log.Println("INSERT: Name: " + name + " | City: " + city)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

