package models

import (
	"backend/function"
	
)

type LoginResponse struct {
	Response  function.User `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}