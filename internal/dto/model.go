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
	Status     string    `db:"status"`
	Accrual    *float64  `db:"accrual" `
	UploadedAt time.Time `db:"uploaded_at"`
	UserID     string    `db:"user_id" `
}

//easyjson:json
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
	UserID      string    `db:"user_id" `
}

type WithdrawalList []struct {
	Number      string    `db:"number" json:"order" `
	Sum         float64   `db:"sum" json:"sum" `
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}

//easyjson:json
type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

//easyjson:json
type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
