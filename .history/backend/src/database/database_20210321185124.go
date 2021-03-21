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

type File struct {
	Sender string
	Filename string
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
    err := model.db.QueryRow("SELECT id FROM users where username = ? AND password=?", username, password).Scan(&id)
	if err != nil {
		return 0, err
	}
    fmt.Println("Id being get from db: ", id)
	return id, nil
}

func (model *DBModel) GetFiles(recipient int) ([]string, error) {
 var filenames []string
 var filename string
 rows, err := model.db.Query("SELECT filename FROM files where recipient_id = ?", recipient)
  if err != nil {
    // handle this error better than this
    return nil,err
  }
  defer rows.Close()
  
  for rows.Next() {
    err = rows.Scan(&filename)
    if err != nil {
      // handle this error
      return nil,err
    }
    fmt.Println(filename)
    filenames = append(filenames,filename)
  }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    return nil,err
  }
  return filenames,err
}
// SELECT files.filename, users.username FROM files INNER JOIN users ON files.sender_id = users.id where files.recipient_id = 6

func (model *DBModel) GetFilesAndSender(recipient int) ([]File, error) {
	var records []File
	var filename string
	var sender string
	rows, err := model.db.Query("SELECT files.filename, users.username FROM files INNER JOIN users ON files.sender_id = users.id where files.recipient_id = ?", recipient)
	 if err != nil {
	   // handle this error better than this
	   return nil,err
	 }
	 defer rows.Close()
	 
	 for rows.Next() {
	   err = rows.Scan(&filename.&sender)
	   if err != nil {
		 // handle this error
		 return nil,err
	   }
	   fmt.Println(filename)
	   filenames = append(filenames,filename)
	 }
	 // get any error encountered during iteration
	 err = rows.Err()
	 if err != nil {
	   return nil,err
	 }
	 return filenames,err
   }

func (model *DBModel) InsertFile(f string,sender int,recipient int) error {
	query := fmt.Sprintf("INSERT INTO files (filename, sender_id, recipient_id) VALUES ('%s', %d, %d)", f, sender, recipient)
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
