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

func (model *DBModel) GetUser(id int) (string, string, error) {
	query :=  fmt.Sprintf("SELECT `username`, `password` FROM `users` WHERE id = $1")
	rows, err := model.db.Query(query, id)
	if err != nil {
		return "", "", err
	}
	var username, password string
	for rows.Next(){
		err := rows.Scan(&username, &password)
		if err != nil {
			return "", "", err
		}
	}
	fmt.Println("Name: ", username, "Password: ", password)
	return username, password, nil
}

func (model *DBModel) InsertUser(p models.Credentials) error {
	query := fmt.Sprintf("INSERT INTO `users` (`username`, `password`) VALUES ($1, $2)")
    _, err := model.db.Query(query, p.Username, p.Password)
    if err != nil {
    	return err
	}
	return nil
}
