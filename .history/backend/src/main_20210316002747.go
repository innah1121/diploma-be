package main

import (
	"backend/function"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Articles []Article

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	// myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/register", signUp).Methods("POST")
	myRouter.HandleFunc("/login", signIn)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var p Credentials
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("username " + p.Username)
	fmt.Println("password " + p.Password)
	user, error := function.InitUser(p.Username, p.Password)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("User is getting registered.Username : " + user.Username)
	fmt.Fprintf(w, "Person: %+v", user)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var p Credentials
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, error := function.GetUser(p.Username, p.Password)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Trying to get user with username : " + user.Username)
	fmt.Fprintf(w, "Person: %+v", p)
	json.NewEncoder(w).Encode(user.Username)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var e Credentials
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
