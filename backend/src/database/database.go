package database

import (
	"backend/models"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int
	Username string
	Password string
}

func Connect() (*sql.DB, error) {
	// Open up our database connection.
	db, err := sql.Open("mysql", "dori:dori@tcp(localhost:3306)/sharesecurely")

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// defer the close after has finished executing
	defer db.Close()

	return db, nil
}

// func Insert(username string, password string) {

//     fmt.Println("INSERT: Username: " + username + " | Password: " + password)
//     insForm, err := db.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
//         if err != nil {
//             panic(err.Error())
//         }
//         insForm.Exec(username, password)
//         fmt.Println("INSERT: Username: " + username + " | Password: " + password)

// }

/*func GetByUsername(username string) {

	selDB, err := db.Query("SELECT * FROM Employee WHERE username=?", username)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(selDB)

}

func GetByUsernameAndPassword(username string, password string) {

	selDB, err := db.Query("SELECT * FROM Employee WHERE username=? AND password=?", username, password)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(selDB)

}*/

func GetUser(db *sql.DB, id string) (string, string, error) {
	query :=  fmt.Sprintf("SELECT name, password FROM users WHERE id = $1")
	rows, err := db.Query(query, id)
	if err != nil {
		return "", "", err
	}
	var name, password string
	for rows.Next(){
		err := rows.Scan(&name, &password)
		if err != nil {
			return "", "", err
		}
	}
	fmt.Println("Name: ", name, "Password: ", password)
	return name, password, nil
}

func InsertUser(db *sql.DB, p models.Credentials) error {
	query := fmt.Sprintf("INSERT INTO users (username, password) VALUES ($1, $2)")
    _, err := db.Query(query, p.Username, p.Password)
    if err != nil {
    	return err
	}
	return nil
}
