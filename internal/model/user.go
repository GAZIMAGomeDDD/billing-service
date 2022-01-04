package model

type User struct {
	ID      string  `json:"id" db:"id"`
	Balance float64 `json:"balance" db:"balance"`
}
