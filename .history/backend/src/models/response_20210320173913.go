package models

import (
	"backend/function"
	
)

type LoginResponse struct {
	User  *function.User `json:"user,omitempty"`
	UserId  int `json:"userId,omitempty"`
	Error error  `json:"error,omitempty"`
}


type FileResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}

type ShareFileResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}
