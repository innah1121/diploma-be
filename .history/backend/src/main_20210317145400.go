package main

import (
	"backend/function"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
	"backend/models"
)



func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/register", signUp).Methods("POST")
	myRouter.HandleFunc("/login", login)
	myRouter.HandleFunc("/storeFile", storeFile).Methods("POST")
	// myRouter.HandleFunc("/appendFile", appendFile).Methods("POST")
	myRouter.HandleFunc("/loadFile", loadFile)
	myRouter.HandleFunc("/shareFile", shareFile)
	
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

// The request is http://localhost:10000/register
// with following body { "username": "dorina", "password":"uka" }
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
	response, _ := json.Marshal(user)
	fmt.Println("User is getting registered.Username : " + user.Username)
	w.Write(response)
}

func login(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
    username := v.Get("username")
	password := v.Get("password")
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.LoginResponse{Response: nil, Error: err})
		w.Write(response)
		return
	}
	user, error := function.GetUser(username, password)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.LoginResponse{Response: nil, Error: err})
		w.Write(response)
		return
	}
	fmt.Println("Trying to get user with username : " + user.Username)
	response, _ := json.Marshal(models.LoginResponse{Response: user, Error: nil})
	w.Write(response)
}
// filename , data  needed
func storeFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    r.ParseMultipartForm(10 << 20)
    // FormFile returns the first file for the given key `myFile`
    // it also returns the FileHeader so we can get the Filename,
    // the Header and the size of the file
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

    // Create a temporary file within our temp-images directory that follows
    // a particular naming pattern
    tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
    if err != nil {
        fmt.Println(err)
    }
    defer tempFile.Close()

    // read all of the contents of our uploaded file into a
    // byte array
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    // write this byte array to our temporary file
    tempFile.Write(fileBytes)
	user, _ := function.GetUser("alice", "fu")
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
//filename  needed
func loadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File loading Endpoint Hit")
	var f models.File
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(f.Filename)
	user, _ := function.GetUser("alice", "fu")
	data, error := user.LoadFile(f.Filename)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	
	fmt.Println(data)
	fmt.Println("Trying to get file with user : " + f.Filename)
	response, _ := json.Marshal(models.ShareFileResponse{Response: "Loaded succesfully", Error: nil})
	w.Write(response)
}
// filename , recipient  needed
func shareFile(w http.ResponseWriter, r *http.Request) {
	var p models.ShareFileRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.ShareFileResponse{Response: "", Error: err})
		w.Write(response)
		return
	}
	user, _ := function.GetUser(p.Username, p.Password)
	data, error := user.ShareFile(p.Filename, p.Recipient)
	
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.ShareFileResponse{Response: "", Error: err})
		w.Write(response)
		return
	}
	fmt.Println(data)
	fmt.Println("Trying to share file with user : " + p.Recipient)
	response, _ := json.Marshal(models.ShareFileResponse{Response: "Shared succesfully", Error: nil})
	w.Write(response)
}


func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
    handleRequests()
}
