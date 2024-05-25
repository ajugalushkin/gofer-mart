package dto

import "time"

//go:generate easyjson -all json.go

//easyjson:json
type Login struct {
	Login    string `json:"login" `
	Password string `json:"password" `
}

//easyjson:json
type User struct {
	ID       string `db:"id"`
	Login    string `db:"login" `
	Password string `db:"password_hash" `
}

//easyjson:json
type Order struct {
	ID         string    `db:"id"`
	Number     string    `db:"number" `
	UploadedAt time.Time `db:"uploaded_at"`
	Status     string    `db:"status"`
	UserID     string    `db:"user_id" `
}

//easyjson:json
type OrderList []Order
