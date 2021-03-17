package models

import (
	"backend/function"
	
)
type RegisterResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}

type LoginResponse struct {
	Response  *function.User `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}