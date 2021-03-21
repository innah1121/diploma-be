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

type DBModel struct {
	db *sql.DB
}

// NewDBModel creates a new database struct
func NewDBModel(database *sql.DB) *DBModel {
	return &DBModel{
		db: database,
	}
}

func Connect() (*sql.DB, error) {
	// Open up our database connection.
	db, err := sql.Open("mysql", "dori:dori@tcp(localhost)/sharesecurely")

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

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

func (model *DBModel) GetUsersPassword(username string) (string, error) {
	var password string
	err := model.db.QueryRow("SELECT password FROM users where username = ?", username).Scan(&password)
	if err != nil {
		return "", err
	}
	fmt.Println("Password getting from db: ", password)
	return password, nil
}

// SELECT * FROM `users` WHERE `username` = "dori" AND `password` = "dori"
func (model *DBModel) GetUserByUsernamePassword(username string,password string) (*User, error) {
	var id int
    var usr string
    var pass string
	err := model.db.QueryRow("SELECT * FROM users where username = ? AND password=?", username, password).Scan(&id,&usr,&pass)
	if err != nil {
		return nil, err
	}
    user := &User{
        Id: id,
        Username: usr,
        Password: pass,
    }
	fmt.Println("User being get from db: ", user)
	return user, nil
}

func (model *DBModel) GetIdByUsernamePassword(username string,password string) (int, error) {
	var id int
    err := model.db.QueryRow("SELECT * FROM users where username = ? AND password=?", username, password).Scan(&id)
	if err != nil {
		return 0, err
	}
    fmt.Println("Id being get from db: ", id)
	return id, nil
}

func (model *DBModel) InsertFile(p models.Credentials) error {
	query := fmt.Sprintf("INSERT INTO users (username, password) VALUES ('%s', '%s')", p.Username, p.Password)
    _, err := model.db.Query(query)
    if err != nil {
    	return err
	}
	return nil
}

func (model *DBModel) InsertUser(p models.Credentials) error {
	query := fmt.Sprintf("INSERT INTO users (username, password) VALUES ('%s', '%s')", p.Username, p.Password)
    _, err := model.db.Query(query)
    if err != nil {
    	return err
	}
	return nil
}
