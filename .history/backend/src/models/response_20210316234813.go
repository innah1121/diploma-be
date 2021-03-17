package models

type LoginResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}