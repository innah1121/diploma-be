package database

import (
	"database/sql"
	"context"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
	"backend/models"
)
var db *sql.DB

func Connect() (*sql.DB, error) {
	// Open up our database connection.
	db, err := sql.Open("mysql", "dori:dori@tcp(localhost:3306)/sharesecurely")

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
	
    fmt.Println("INSERT: Username: " + username + " | Password: " + password)
    insForm, err := db.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
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
func insert(db *sql.DB, p models.Credentials) error {  
    query := "INSERT INTO users(username, password) VALUES (?, ?)"
    ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    stmt, err := db.PrepareContext(ctx, query)
    if err != nil {
        fmt.Printf("Error %s when preparing SQL statement", err)
        return err
    }
    defer stmt.Close()
    res, err := stmt.ExecContext(ctx, p.Username, p.Password)
    if err != nil {
        fmt.Printf("Error %s when inserting row into products table", err)
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        fmt.Printf("Error %s when finding rows affected", err)
        return err
    }
    fmt.Printf("%d products created ", rows)
    return nil
}

