package main

import (
	"backend/database"
	"backend/function"
	"backend/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var dbModel *database.DBModel

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/register", signUp).Methods("POST")
	myRouter.HandleFunc("/login", login)
	myRouter.HandleFunc("/storeFile", storeFile).Methods("POST")
	myRouter.HandleFunc("/loadFile", loadFile)
	myRouter.HandleFunc("/shareFile", shareFile)
	log.Fatal(http.ListenAndServe(":10000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(myRouter)))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var p models.Credentials
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(err.Error())
		w.Write(response)
		return
	}
	fmt.Println("username " + p.Username)
	fmt.Println("password " + p.Password)
	user, error := function.InitUser(p.Username, p.Password)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal("User already registered")
		w.Write(response)
		return
	}
	dbModel.InsertUser(p)
	username, _, _ := dbModel.GetUser(1)
	fmt.Println(username + " taken from db")
	response, _ := json.Marshal(user)
	fmt.Println("User is getting registered.Username : " + user.Username)
	w.Write(response)
}

func login(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
    username := v.Get("username")
	password := v.Get("password")
	// check if exists in my local db ,if return result i go on with the other logic
	function.InitUser(username, password)
	user, error := function.GetUser(username, password)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.LoginResponse{Response: nil, Error: error})
		w.Write(response)
		return
	}
	fmt.Println("Trying to get user with username : " + user.Username)
	response, _ := json.Marshal(models.LoginResponse{Response: user, Error: nil})
	w.Write(response)
}

func storeFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	v := r.URL.Query()
    username := v.Get("username")
	password := v.Get("password")
	function.InitUser(username, password)
    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    r.ParseMultipartForm(10 << 20)
    // FormFile returns the first file for the given key `myFile`
    // it also returns the FileHeader so we can get the Filename,
    // the Header and the size of the file
    file, handler, err := r.FormFile("filename")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    // Create a temporary file within our temp-images directory that follows
    // a particular naming pattern
    // tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
    // if err != nil {
    //     fmt.Println(err)
    // }
    // defer tempFile.Close()

    // read all of the contents of our uploaded file into a
    // byte array
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    // write this byte array to our temporary file
    //tempFile.Write(fileBytes)
	user, _ := function.GetUser(username, password)
	fmt.Println(fileBytes)
	user.StoreFile(handler.Filename, fileBytes)
	// (*function.User)(nil).StoreFile(handler.Filename, fileBytes)
    // return that we have successfully uploaded our file!
    fmt.Fprintf(w, "Successfully Uploaded File\n")
}
// filename , data  needed
func appendFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File append Endpoint Hit")
    r.ParseMultipartForm(10 << 20)
    file, handler, err := r.FormFile("myFile")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)
    tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
    if err != nil {
        fmt.Println(err)
    }
    defer tempFile.Close()
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    tempFile.Write(fileBytes)
	user, _ := function.GetUser("alice", "fubar")
	user.AppendFile(handler.Filename, fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func loadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File loading Endpoint Hit")
	v := r.URL.Query()
    filename := v.Get("filename")
	username := v.Get("username")
	password := v.Get("password")
	// mbase dhe init user
	user, _ := function.GetUser(username, password)
	fmt.Println(user)
	data, err := user.LoadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	fmt.Println(data)
	fmt.Println("Trying to get file with name : " + filename)
	response, _ := json.Marshal(models.ShareFileResponse{Response: "Loaded succesfully", Error: nil})
	w.Write(response)
}
// filename , recipient  needed
func shareFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File sharing Endpoint Hit")
	v := r.URL.Query()
    filename := v.Get("filename")
	username := v.Get("username")
	password := v.Get("password")
	recipient := v.Get("recipient")
	/*if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.ShareFileResponse{Response: "", Error: err})
		w.Write(response)
		return
	}*/
	
	user, _ := function.GetUser(username, password)
	fmt.Println("i got usr maybe")
	user.LoadFile(filename)
	fmt.Println("file loading passed")
	function.InitUser(recipient, password)
	data, error := user.ShareFile(filename, recipient)
	fmt.Println("i might be stuck in share")
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.ShareFileResponse{Response: "", Error: error})
		w.Write(response)
		return
	}
	fmt.Println(data)
	fmt.Println("Trying to share file with user : " + recipient)
	response, _ := json.Marshal(models.ShareFileResponse{Response: "Shared succesfully", Error: nil})
	w.Write(response)
}

func recieveFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File recieving Endpoint Hit")
	
	var p models.RecieveFileRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	user, _ := function.GetUser(p.SenderUsr, p.SenderPass)
	user2, _ := function.GetUser(p.RecipientUsr, p.RecipientPass)
	magic_string, er := user.ShareFile(p.Filename, p.RecipientUsr)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}
	error := user2.ReceiveFile("file2", p.SenderUsr, magic_string)
	
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Trying to share file with user : " + p.RecipientUsr)
	response, _ := json.Marshal(models.ShareFileResponse{Response: "Recieved succesfully", Error: nil})
	w.Write(response)
}


func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	db, err := database.Connect()
	dbModel = database.NewDBModel(db)
	if err != nil {
		fmt.Println("Error connecting db")
		os.Exit(1)
	}
	fmt.Println("Connected to db.")
	handleRequests()
}
