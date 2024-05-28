package dto

import "time"

//go:generate easyjson -all json.go

//easyjson:json
type User struct {
	Login    string `json:"login" db:"login" `
	Password string `json:"password" db:"password" `
}

type Order struct {
	Number     string    `db:"number" `
	UploadedAt time.Time `db:"uploaded_at"`
	Status     string    `db:"status"`
	Accrual    float64   `db:"accrual" `
	UserID     string    `db:"user_id" `
}

type OrderList []Order

//easyjson:json
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	Number      string    `db:"number" `
	Sum         float64   `db:"sum" `
	ProcessedAt time.Time `db:"processed_at"`
}
type WithdrawalList []Withdrawal

//easyjson:json
type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
