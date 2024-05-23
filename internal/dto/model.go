package dto

//go:generate easyjson -all json.go

//easyjson:json
type Login struct {
	Login    string `json:"login" `
	Password string `json:"password" `
}

//easyjson:json
type User struct {
	ID       string `db:"user_id"`
	Login    string `db:"login" `
	Password string `db:"password_hash" `
}

//easyjson:json
type Order struct {
	ID          string `db:"order_id"`
	OrderNumber string `db:"order_number" `
	UserID      string `db:"user_id" `
}
