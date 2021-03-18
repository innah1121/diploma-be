package models

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type File struct {
	Filename string `json:"filename"`
}

type ShareFileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Filename string `json:"filename"`
	Recipient string `json:"recipient"`
}

type RecieveFileRequest struct {
	Filename string `json:"filename"`
	SenderUsr string `json:"senderU"`
	SenderPass string `json:"senderP"`
	MagicString string `json:"magic_string"`
}

type RevokeFileRequest struct {
	Filename string `json:"filename"`
	TargetUsername string `json:"target_username"`
}