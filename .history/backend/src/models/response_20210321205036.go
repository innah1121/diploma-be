package models

import (
	"backend/function"
	
)
var dbModel *database.DBModel

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
	Files  []dbModel.File `json:"files,omitempty"`
	Error error  `json:"error,omitempty"`
}

type ShareFileResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}
