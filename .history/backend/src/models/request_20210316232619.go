package models

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type File struct {
	Filename string `json:"filename"`
}

type SharedFile struct {
	Filename string `json:"filename"`
	Recipient string `json:"recipient"`
}
