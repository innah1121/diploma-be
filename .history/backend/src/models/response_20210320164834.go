package models

import (
	"backend/function"
	
)

type LoginResponse struct {
	Response  *function.User `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}
type LoginResponse2 struct {
	Response  int `json:"response,omitempty"`
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
