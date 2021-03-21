package models

import (
	"backend/function"
	
	
)


type FileDb struct {
	Sender string
	Filename string
}

type LoginResponse struct {
	User  *function.User `json:"user,omitempty"`
	UserId  int `json:"userId,omitempty"`
	Error error  `json:"error,omitempty"`
}


type FileResponse struct {
	Files  []string `json:"files,omitempty"`
	Error error  `json:"error,omitempty"`
}

type FileResponse2 struct {
	Files  []FileDb `json:"files,omitempty"`
	Error error  `json:"error,omitempty"`
}

type ShareFileResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}
