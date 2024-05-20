package dto

//go:generate easyjson -all json.go

//easyjson:json
type Login struct {
	Login    string `json:"login" `
	Password string `json:"password" `
}

//easyjson:json
type User struct {
	Id       string `db:"id"`
	Login    string `db:"user_id" `
	Password string `db:"password_hash" `
}
