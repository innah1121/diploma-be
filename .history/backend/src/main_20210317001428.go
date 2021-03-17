package main

import (
	"backend/function"
	"encoding/json"
	"errors"
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
	myRouter.HandleFunc("/appendFile", appendFile).Methods("POST")
	myRouter.HandleFunc("/loadFile", loadFile)
	
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
	var p models.Credentials
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.LoginResponse{Response: nil, Error: err})
		w.Write(response)
		return
	}
	user, error := function.GetUser(p.Username, p.Password)
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


// The request is http://localhost:10000/login/username/dorina
// with following body { "password":"uka" }
func signIn(w http.ResponseWriter, r *http.Request) {
	var p models.Credentials
	vars := mux.Vars(r)
	p.Username = vars["username"] // you need to specify "username" string in the url
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.RegisterResponse{Response: "", Error: err})
		w.Write(response)
		return
	}
	user, error := function.GetUser(p.Username, p.Password)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(models.RegisterResponse{Response: "User not found", Error: err})
		w.Write(response)
		return
	}
	fmt.Println("Trying to get user with username : " + user.Username)
	response, _ := json.Marshal(models.RegisterResponse{Response: "Login successfully", Error: nil})
	w.Write(response)
}

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
	user, _ := function.GetUser("alice", "fubar")
	user.StoreFile(handler.Filename, fileBytes)
	// (*function.User)(nil).StoreFile(handler.Filename, fileBytes)
    // return that we have successfully uploaded our file!
    fmt.Fprintf(w, "Successfully Uploaded File\n")
}

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
	var f models.File
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(f.Filename)
	user, _ := function.GetUser("alice", "fubar")
	data, error := user.LoadFile(f.Filename)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Trying to get file with name : " + f.Filename)
	fmt.Fprintf(w, "file: %+v", data)
	json.NewEncoder(w).Encode(data)
}

func shareFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File loading Endpoint Hit")
	var f models.SharedFile
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(f.Filename)
	fmt.Println(f.Recipient)
	user, _ := function.GetUser("alice", "fubar")
	data, error := user.LoadFile(f.Filename)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Trying to get file with name : " + f.Filename)
	fmt.Fprintf(w, "file: %+v", data)
	json.NewEncoder(w).Encode(data)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var e models.Credentials
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	errorResponse(w, "Success", http.StatusOK)
	return
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
    handleRequests()
}
