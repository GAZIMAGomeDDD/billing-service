package model

import "time"

type Transaction struct {
	ID        string    `json:"id" db:"id"`
	ToID      *string   `json:"to_id,omitempty" db:"to_id"`
	FromID    *string   `json:"from_id,omitempty" db:"from_id"`
	Money     float64   `json:"money" db:"money"`
	Method    string    `json:"method" db:"method"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ListOfTransactionsQuery struct {
	UserID string `json:"user_id" example:"b91a95a4-078f-4afd-b11c-4850eb65e784"`
	Limit  int    `json:"limit" example:"5"`
}

type GetTransaction struct {
	ID string `json:"id" example:"b91a95a4-078f-4afd-b11c-4850eb65e784"`
}
