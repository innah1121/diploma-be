package models

type RegisterResponse struct {
	User  string `json:"user,omitempty"`
	Error error  `json:"error,omitempty"`
}

type LoginResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}